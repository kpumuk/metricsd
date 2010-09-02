package main

import (
    "bytes"
    "flag"
    "net"
    "os"
    "os/signal"
    "path"
    "strconv"
    "strings"
    "time"
    "./config"
    "./logger"
    "./types"
    "./writers"
)

var (
    log logger.Logger
    /* Slices */
    slices *types.Slices
)

func lookupHost(addr *net.UDPAddr) string {
    return addr.IP.String()
}

func process(addr *net.UDPAddr, buf string, msgchan chan<- *types.Message) {
    log.Debug("Processing message from %s: %s", addr, buf)
    var source, name, svalue string
    var fields []string

    // Multiple metrics in a single message
    for _, msg := range strings.Split(buf, ";", -1) {
        // Check if the message contains a source name
        if idx := strings.Index(msg, "@"); idx >= 0 {
            source = msg[0:idx]
            msg = msg[idx+1:]
        } else {
            source = lookupHost(addr)
        }

        // Retrieve the metric name
        fields = strings.Split(msg, ":", -1)
        if len(fields) < 2 {
            log.Debug("Message format is not valid: %s", buf)
            return
        } else {
            name = fields[0]
            svalue = fields[1]
        }

        // Parse the value
        if value, error := strconv.Atoi(svalue); error != nil {
            log.Debug("Number %s is not valid: %s", svalue, error)
        } else {
            msgchan <- types.NewMessage(source, name, value)
        }
    }
}

func listen(msgchan chan<- *types.Message, quit chan bool) {
    log.Debug("Starting listener on %s", config.GlobalConfig.UDPAddress)

    // Listen for requests
    listener, error := net.ListenUDP("udp", config.GlobalConfig.UDPAddress)
    if error != nil {
        log.Fatal("Cannot listen: %s", error)
        os.Exit(1)
    }
    // Ensure listener will be closed on return
    defer listener.Close()

    // Timeout is 0.1 second
    listener.SetTimeout(100000000)
    listener.SetReadTimeout(100000000)

    message := make([]byte, 256)
    for {
        if _, ok := <-quit; ok {
            log.Debug("Shutting down listener...")
            return
        }

        n, addr, error := listener.ReadFromUDP(message)
        if error != nil {
            if addr != nil {
                log.Debug("Cannot read UDP from %s: %s\n", addr, error)
            }
            continue
        }
        buf := bytes.NewBuffer(message[0:n])
        process(addr, buf.String(), msgchan)
    }
}

func msgSlicer(msgchan <-chan *types.Message) {
    for {
        message := <-msgchan
        slices.Add(message)
    }
}

func initialize() {
    // Initialize options parser
    var slice, write, debug int
    var listen, data, cfg string
    var test, batch bool
    flag.StringVar(&cfg,     "config",  config.DEFAULT_CONFIG_PATH,    "Set the path to config file")
    flag.StringVar(&listen,  "listen",  config.DEFAULT_LISTEN,         "Set the port (+optional address) to listen at")
    flag.StringVar(&data,    "data",    config.DEFAULT_DATA_DIR,       "Set the data directory")
    flag.IntVar   (&debug,   "debug",   int(config.DEFAULT_SEVERITY),  "Set the debug level, the lower - the more verbose (0-5)")
    flag.IntVar   (&slice,   "slice",   config.DEFAULT_SLICE_INTERVAL, "Set the slice interval in seconds")
    flag.IntVar   (&write,   "write",   config.DEFAULT_WRITE_INTERVAL, "Set the write interval in seconds")
    flag.BoolVar  (&batch,   "batch",   config.DEFAULT_BATCH_WRITES,   "Set the value indicating whether batch RRD updates should be used")
    flag.BoolVar  (&test,    "test",    false,                         "Validate config file and exit")
    flag.Parse()

    // Load config from a config file
    config.GlobalConfig.Load(cfg)
    if test { os.Exit(0) }

    // Override options with values passed in command line arguments
    // (but only if they have a value different from a default one)
    if listen != config.DEFAULT_LISTEN {
        config.GlobalConfig.Listen        = listen
    }
    if data != config.DEFAULT_DATA_DIR {
        config.GlobalConfig.DataDir       = data
    }
    if debug != int(config.DEFAULT_SEVERITY) {
        config.GlobalConfig.LogLevel      = debug
    }
    if slice != config.DEFAULT_SLICE_INTERVAL {
        config.GlobalConfig.SliceInterval = slice
    }
    if write != config.DEFAULT_WRITE_INTERVAL {
        config.GlobalConfig.WriteInterval = write
    }
    if batch != config.DEFAULT_BATCH_WRITES {
        config.GlobalConfig.BatchWrites   = batch
    }

    // Make data dir path absolute
    if len(config.GlobalConfig.DataDir) == 0 || config.GlobalConfig.DataDir[0] != '/' {
        wd, _ := os.Getwd()
        config.GlobalConfig.DataDir = path.Join(wd, config.GlobalConfig.DataDir)
    }

    // Create logger
    config.GlobalConfig.Logger = logger.NewConsoleLogger(logger.Severity(config.GlobalConfig.LogLevel))
    log = config.GlobalConfig.Logger
    log.Debug("%s", config.GlobalConfig)

    // Ensure data directory exists
    if _, err := os.Stat(data); err != nil {
        os.MkdirAll(data, 0755)
    }

    // Resolve listen address
    address, error := net.ResolveUDPAddr(config.GlobalConfig.Listen)
    if error != nil {
        log.Fatal("Cannot parse \"%s\": %s", config.GlobalConfig.Listen, error)
        os.Exit(1)
    }
    config.GlobalConfig.UDPAddress = address

    // Initialize slices structure
    slices = types.NewSlices(config.GlobalConfig.SliceInterval)
}

func rollupSlices(active_writers []writers.Writer, force bool) {
    log.Debug("Rolling up slices")

    if config.GlobalConfig.BatchWrites {
        closedSampleSets := slices.ExtractClosedSampleSets(force)
        for _, writer := range active_writers {
            writers.BatchRollup(writer, closedSampleSets)
        }
    } else {
        closedSlices := slices.ExtractClosedSlices(force)
        closedSlices.Do(func(elem interface {}) {
            slice := elem.(*types.Slice)
            for _, set := range slice.Sets {
                for _, writer := range active_writers {
                    writers.Rollup(writer, set)
                }
            }
        })
    }
}

func dumper(active_writers []writers.Writer, quit chan bool) {
    ticker := time.NewTicker(int64(config.GlobalConfig.WriteInterval) * 1000000000)
    defer ticker.Stop()

    for {
        if _, ok := <-quit; ok {
            log.Debug("Shutting down dumper...")
            return
        }

        <-ticker.C
        rollupSlices(active_writers, false)
    }
}

func main() {
    initialize()

    // Quit channel. Should be blocking (non-bufferred), so sender
    // will wait till receiver will accept message (and shut down)
    quit := make(chan bool)

    // Messages channel
    msgchan := make(chan *types.Message, 1000)
    go msgSlicer(msgchan)

    active_writers := []writers.Writer {
        &writers.Quartiles {},
        &writers.YesOrNo   {},
    }

    go listen(msgchan, quit)
    go dumper(active_writers, quit)

    for sig := range signal.Incoming {
        var usig = sig.(signal.UnixSignal)
        if usig == 1 || usig == 2 || usig == 15 {
            log.Warn("Received signal: %s", sig)
            if usig == 2 || usig == 15 {
                log.Warn("Shutting down everything...")
                // We have two background processes, so wait for both
                quit <- true
                quit <- true
            }
            rollupSlices(active_writers, true)
            if usig == 2 || usig == 15 { return }
        }
    }
}
