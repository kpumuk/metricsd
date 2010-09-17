package main

import (
    "flag"
    "fmt"
    "io"
    "log"
    "net"
    "rand"
    "runtime"
    "strings"
    "time"
)

func send(address *net.UDPAddr, source, key string, value int) {
    conn, error := net.DialUDP("udp4", nil, address)
    if error != nil {
        log.Stderrf("Failed to connect to %s", address)
    }
    defer conn.Close()

    data := fmt.Sprintf("%s@%s:%d", source, key, value)
    io.Copy(conn, strings.NewReader(data))
}

func main() {
    var address, source, key string
    var count, step, threads int
    var delay int64
    var sourcecnt, keycnt int
    flag.StringVar(&address,   "address",   "127.0.0.1:6311", "Set the port (+optional address) to send packets to")
    flag.StringVar(&source,    "source",    "app%d",      "Set the source name (pakets data will be \"source@key:idx\")")
    flag.IntVar   (&sourcecnt, "sourcecnt", 10,           "Set the number of sources to send from (when \"key%d\" substitution pattern in -source is used)")
    flag.StringVar(&key,       "key",       "profile_%d", "Set the key name (pakets data will be \"key:idx\")")
    flag.IntVar   (&keycnt,    "keycnt",    100,          "Set the number of metrics (when \"metric%d\" substitution pattern in -key is used)")
    flag.IntVar   (&count,     "count",     1000,         "Set the number of packets to send")
    flag.Int64Var (&delay,     "delay",     1000000,      "Set the delay between packets in nanoseconds (10^-9)")
    flag.IntVar   (&step,      "step",      100,          "Log step (how many packets to send between logging)")
    flag.IntVar   (&threads,   "threads",   10,           "Set the number of active threads")
    flag.Parse()

    udp_address, error := net.ResolveUDPAddr(address)
    if error != nil {
        log.Exitf("Cannot parse \"%s\": %s", address, error)
    }

    runtime.GOMAXPROCS(threads + 1)
    rand.Seed(time.Nanoseconds())

    tasks := make(chan int, threads)
    for i := 1; i <= threads; i++ {
        go func(idx int) {
            ticker := time.NewTicker(delay)
            defer ticker.Stop()

            log.Stdoutf("Started thread #%d", idx)
            for {
                <-ticker.C
                task := <-tasks
                send(udp_address, fmt.Sprintf(source, rand.Intn(sourcecnt)), fmt.Sprintf(key, rand.Intn(keycnt)), task % step)
            }
        }(i)
    }

    for sent := 1; sent <= count; sent++ {
        tasks <- sent
        if sent % step == 0 { log.Stdoutf("Processed %d packets of %d", sent, count) }
    }
    time.Sleep(delay * 2)
}
