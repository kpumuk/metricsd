package main

import (
    "bytes"
    "flag"
    "fmt"
    "log"
    "net"
    "time"
)

func send(address *net.UDPAddr, key string, value int) {
    conn, error := net.DialUDP("udp", nil, address)
    if error != nil {
        log.Stderrf("Failed to connect to %s", address)
    }
    defer conn.Close()

    data := fmt.Sprintf("%s:%d", key, value)
    buf := bytes.NewBufferString(data)

    conn.Write(buf.Bytes())
}

func main() {
    var address, key string
    var count, delay, step int
    flag.StringVar(&address, "address", "127.0.0.1:6311", "Set the port (+optional address) to send packets to")
    flag.StringVar(&key,     "key",     "profile", "Set the key name (pakets data will be \"key:idx\")")
    flag.IntVar   (&count,   "count",   1000,      "Set the number of packets to send")
    flag.IntVar   (&delay,   "delay",   1000000,   "Set the delay between packets in nanoseconds (10^-9)")
    flag.IntVar   (&step,    "step",    100,       "Log step (how many packets to send between logging)")
    flag.Parse()

    udp_address, error := net.ResolveUDPAddr(address)
    if error != nil {
        log.Exitf("Cannot parse \"%s\": %s", address, error)
    }

    log.Stdout(udp_address)

    ticker := time.NewTicker(int64(delay))
    defer ticker.Stop()

    for sent := 1; sent <= count; sent++ {
        <-ticker.C;
        send(udp_address, key, sent)
        if sent % step == 0 { log.Stdoutf("Processed %d packets of %d", sent, count) }
    }
}
