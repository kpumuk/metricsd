package main

import (
	"net"
	"os"
	"os/signal"
	"runtime"
	"time"
	"metricsd/config"
	"metricsd/logger"
	"metricsd/parser"
	"metricsd/writers"
	"metricsd/stdlib"
	"metricsd/types"
	"metricsd/web"
)

var (
	log                 logger.Logger     /* Logger instance */
	hostLookupCache     map[string]string /* DNS names cache */
	timeline            *types.Timeline   /* Timeline */
	eventsReceived      int64             /* Events received */
	totalEventsReceived int64             /* Total Events received */
	bytesReceived       int64             /* Bytes sent */
	totalBytesReceived  int64             /* Total bytes sent */
	activeWriters       []writers.Writer  /* The list of active writers */
)

const (
	runningProcesses = 3
)

func main() {
	// Initialize MetricsD
	initialize()

	// Quit channel. Should be blocking (non-bufferred), so sender
	// will wait until receiver accepts the message
	// (and then will shut himself down).
	quit := make(chan bool)

	// Active writers
	activeWriters = []writers.Writer{
		&writers.Count{},
		&writers.Quartiles{},
		&writers.Percentiles{},
	}

	// Start background Go routines
	go listen(quit)
	go stats(quit)
	go dumper(activeWriters, quit)
	go web.Start()

	// Handle signals
	handleSignals(quit)
}

func initialize() {
	// Initialize options parser
	parseCommandLineArguments()

	// Create logger
	config.Logger = logger.NewConsoleLogger(logger.Severity(config.LogLevel))
	log = config.Logger
	log.Debug("%s", config.String())

	// Ensure data directory exists
	if _, err := os.Stat(config.DataDir); err != nil {
		os.MkdirAll(config.DataDir, 0755)
	}

	// Resolve listen address
	address, error := net.ResolveUDPAddr("udp", config.Listen)
	if error != nil {
		log.Fatal("Cannot parse \"%s\": %s", config.Listen, error)
		os.Exit(1)
	}
	config.UDPAddress = address

	// Initialize slices structure
	timeline = types.NewTimeline(config.SliceInterval)

	// Initialize host lookup cache
	if config.LookupDns {
		hostLookupCache = make(map[string]string)
	}

	// Disable memory profiling to prevent panics reporting
	runtime.MemProfileRate = 0
}

func handleSignals(quit chan<- bool) {
	for sig := range signal.Incoming {
		var usig = sig.(os.UnixSignal)
		if usig == os.SIGHUP || usig == os.SIGINT || usig == os.SIGTERM {
			log.Warn("Received signal: %s", sig)
			if usig == os.SIGINT || usig == os.SIGTERM {
				log.Warn("Shutting down everything...")
				// We have several background processes, so wait for all of them
				for i := 1; i <= runningProcesses; i++ {
					log.Debug("... waiting for process %d of %d", i, runningProcesses)
					quit <- true
				}
				log.Warn("... done!")
			}
			rollupSlices(activeWriters, true)
			if usig == os.SIGINT || usig == os.SIGTERM {
				return
			}
		}
	}
}

/***** Go routines ************************************************************/

func listen(quit <-chan bool) {
	log.Debug("Starting listener on %s", config.UDPAddress)

	// Listen for requests
	listener, error := net.ListenUDP("udp", config.UDPAddress)
	if error != nil {
		log.Fatal("Cannot listen: %s", error)
		os.Exit(1)
	}
	// Ensure listener will be closed on return
	defer listener.Close()

	// Timeout is 0.1 second
	listener.SetTimeout(1e8)
	listener.SetReadTimeout(1e8)

	data := make([]byte, 256)
	for {
		select {
		case <-quit:
			log.Debug("Shutting down listener...")
			return
		default:
			n, addr, error := listener.ReadFromUDP(data)
			if error != nil {
				if addr != nil {
					log.Debug("Cannot read UDP from %s: %s\n", addr, error)
				}
				continue
			}
			process(addr, string(data[0:n]))
		}
	}
}

func stats(quit <-chan bool) {
	ticker := time.NewTicker(1e9)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			log.Debug("Shutting down stats...")
			return
		case <-ticker.C:
			timeline.Add(types.NewEvent("all", "metricsd.events.count", int(eventsReceived)))
			timeline.Add(types.NewEvent("all", "metricsd.traffic_in", int(bytesReceived)))
			timeline.Add(types.NewEvent("all", "metricsd.memory.used", int(runtime.MemStats.Alloc/1024)))
			timeline.Add(types.NewEvent("all", "metricsd.memory.system", int(runtime.MemStats.Sys/1024)))

			log.Debug("Processed %d events (%d bytes)", eventsReceived, bytesReceived)

			eventsReceived = 0
			bytesReceived = 0
		}
	}
}

func dumper(activeWriters []writers.Writer, quit <-chan bool) {
	ticker := time.NewTicker(int64(config.WriteInterval) * 1e9)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			log.Debug("Shutting down dumper...")
			return
		case <-ticker.C:
			rollupSlices(activeWriters, false)
		}
	}
}

/***** Helper functions *******************************************************/

func process(addr *net.UDPAddr, buf string) {
	bytesReceived += int64(len(buf))
	totalBytesReceived += int64(len(buf))
	parser.Parse(buf, func(event *types.Event, err os.Error) {
		if err == nil {
			if event.Source == "" {
				event.Source = lookupHost(addr)
			}
			timeline.Add(event)
			eventsReceived++
			totalEventsReceived++
		} else {
			log.Debug("Error while parsing an event: %s", err)
		}
	})
}

func lookupHost(addr *net.UDPAddr) (hostname string) {
	ip := addr.IP.String()
	if !config.LookupDns {
		return ip
	}

	// Do we have resolved this address before?
	if _, found := hostLookupCache[ip]; found {
		return hostLookupCache[ip]
	}

	// Try to lookup
	hostname, error := stdlib.GetRemoteHostName(ip)
	if error != nil {
		log.Debug("Error while resolving host name %s: %s", addr, error)
		return ip
	}
	// Cache the lookup result
	hostLookupCache[ip] = hostname

	return
}

func rollupSlices(activeWriters []writers.Writer, force bool) {
	log.Debug("Rolling up timeline")
	startTime := time.Nanoseconds()

	if config.BatchWrites {
		closedSampleSets := timeline.ExtractClosedSampleSets(force)
		for _, writer := range activeWriters {
			writers.BatchRollup(writer, closedSampleSets)
		}
	} else {
		closedSlices := timeline.ExtractClosedSlices(force)
		for _, slice := range closedSlices {
			for _, set := range slice.Sets {
				for _, writer := range activeWriters {
					writers.Rollup(writer, set)
				}
			}
		}
	}
	log.Debug("... timeline rolled up, took %v seconds", float64(time.Nanoseconds()-startTime)/1e9)
}
