package main

import (
    "bytes"
    "flag"
    "log"
    "net"
    "os"
    "path"
    "strconv"
    "strings"
    "time"
    "./types"
    "./writers"
)

var (
    /**
     * Configuration, indexed by keywords, for instance "debug" or "servername"
     */
    globalConfig map[string] interface{}
    /**
     * Not mandatory but it is simpler to use than
     * globalConfig["debug"], which has type interface{}
     */
    debug int
    /* DNS names cache */
    hostLookupCache map[string] string
    /* Slices */
    slices *types.Slices
)

func lookupHost(addr *net.UDPAddr) string {
    ip := addr.IP.String()
    if _, found := hostLookupCache[ip]; found { return hostLookupCache[ip] }

    cname, _, error := net.LookupHost(ip)
    if error != nil {
        // if debug > 1 { log.Stderrf("Host lookup failed for IP %s: %s", ip, error) }
        return ip
    }
    hostLookupCache[ip] = cname
    return cname
}

func process(addr *net.UDPAddr, buf string, msgchan chan<- *types.Message) {
    if debug > 2 { log.Stdoutf("Processing message from %s: %s", addr, buf) }

    fields := strings.Split(buf, ":", 2)

    if value, error := strconv.Atoi(fields[1]); error != nil {
        if debug > 1 { log.Stderrf("Number %s is not valid: %s", fields[1], error) }
    } else {
        msgchan <- types.NewMessage(lookupHost(addr), fields[0], value)
    }
}

func listen(msgchan chan<- *types.Message) {
    if debug > 2 { log.Stdoutf("Starting listener on %s", globalConfig["address"]) }

    // Listen for requests
    listener, error := net.ListenUDP("udp", globalConfig["address"].(*net.UDPAddr))
    if error != nil {
        log.Exitf("Cannot listen: %s", error)
    }
    // Ensure listener will be closed on return
    defer listener.Close()

    message := make([]byte, 256)
    for {
        n, addr, error := listener.ReadFromUDP(message)
        if error != nil {
            if debug > 1 { log.Stderrf("Cannot read UDP from %s: %s\n", addr, error) }
            continue
        }
        buf := bytes.NewBuffer(message[0:n])
        process(addr, buf.String(), msgchan)
    }
}

func msgSlicer(msgchan <-chan *types.Message) {
    for {
        message := <-msgchan
        log.Stdoutf("Slicing message: %s", message)
        slices.Add(message)
    }
}

func initialize() {
    var slice, write int
    var listen, data string
    flag.StringVar(&listen, "listen", "0.0.0.0:6311", "Set the port (+optional address) to listen at")
    flag.StringVar(&data,   "data",   "", "Set the data directory")
    flag.IntVar   (&debug,  "debug",  0,  "Set the debug level, the higher, the more verbose")
    flag.IntVar   (&slice,  "slice",  10, "Set the slice interval in seconds")
    flag.IntVar   (&write,  "write",  60, "Set the write interval in seconds")
    flag.Parse()

    if len(data) == 0 || data[0] != '/' {
        wd, _ := os.Getwd()
        data = path.Join(wd, data)
    }
    if debug > 2 { log.Stdout("Initializing configuration") }

    globalConfig = make(map[string] interface{})
    globalConfig["listen"] = listen
    globalConfig["data"]   = data
    globalConfig["debug"]  = debug
    globalConfig["slice"]  = slice
    globalConfig["write"]  = write

    if _, err := os.Stat(data); err != nil {
        os.MkdirAll(data, 0755)
    }

    hostLookupCache = make(map[string] string)

    address, error := net.ResolveUDPAddr(listen)
    if error != nil {
        log.Exitf("Cannot parse \"%s\": %s", listen, error)
    }

    globalConfig["address"] = address

    slices = types.NewSlices(&globalConfig)
}

func rollupSlices(active_writers []writers.Writer) {
    if debug > 2 { log.Stdout("Rolling up slices") }

    closedSlices := slices.ExtractClosedSlices(false)
    closedSlices.Do(func(elem interface {}) {
        slice := elem.(*types.Slice)
        for _, set := range slice.Sets {
            for _, writer := range active_writers {
                writer.Rollup(set.Time, set.Key, set.Values)
            }
        }
    })
}

func main() {
    initialize()

    // Messages channel
    msgchan := make(chan *types.Message)
    go msgSlicer(msgchan)

    active_writers := make([]writers.Writer, 2)
    active_writers[0] = &writers.Quartiles { Data: globalConfig["data"].(string) }
    active_writers[1] = &writers.YesOrNo   { Data: globalConfig["data"].(string) }

    ticker := time.NewTicker(int64(globalConfig["write"].(int)) * 1000000000) // 10^9
    defer ticker.Stop()
    go func() {
        for {
            <-ticker.C;
            rollupSlices(active_writers)
        }
    }()

    listen(msgchan)
}
