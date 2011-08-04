package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"rand"
	"runtime/pprof"
	"os"
	"time"
)

func main() {
	var address, source, key, cpuprofile string
	var count, step, threads int
	var delay int64
	var sourcecnt, keycnt int
	flag.StringVar(&address, "address", "127.0.0.1:6311", "Set the port (+optional address) to send packets to")
	flag.StringVar(&source, "source", "app%03d", "Set the source name (pakets data will be \"source@key:idx\")")
	flag.IntVar(&sourcecnt, "sourcecnt", 10, "Set the number of sources to send from (when \"source%d\" substitution pattern in -source is used)")
	flag.StringVar(&key, "key", "benchmark.metric%03d", "Set the key name (pakets data will be \"key:idx\")")
	flag.IntVar(&keycnt, "keycnt", 100, "Set the number of metrics (when \"metric%d\" substitution pattern in -key is used)")
	flag.IntVar(&count, "count", 1000, "Set the number of packets to send")
	flag.Int64Var(&delay, "delay", 1000000, "Set the delay between packets in nanoseconds (10^-9)")
	flag.IntVar(&step, "step", 100, "Log step (how many packets to send between logging)")
	flag.IntVar(&threads, "threads", 10, "Set the number of active threads")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "Write CPU profile to the file")
	flag.Parse()

	udp_address, error := net.ResolveUDPAddr("udp", address)
	if error != nil {
		log.Fatalf("Cannot parse \"%s\": %s", address, error)
	}

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatalf("Error while creating file \"%s\" for CPU profile: %s", cpuprofile, err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	rand.Seed(time.Nanoseconds())

	tasks := make(chan int, threads)
	for i := 1; i <= threads; i++ {
		go func(idx int) {
			ticker := time.NewTicker(delay)
			defer ticker.Stop()

			conn, error := net.DialUDP("udp4", nil, udp_address)
			if error != nil {
				log.Fatalf("Failed to connect to %s", address)
			}

			buf := bytes.NewBuffer(make([]byte, 0, 20))

			log.Printf("Started thread #%d", idx)
			for {
				<-ticker.C
				task := <-tasks

				fmt.Fprintf(buf, source, rand.Intn(sourcecnt))
				buf.WriteRune('@')
				fmt.Fprintf(buf, key, rand.Intn(keycnt))
				fmt.Fprintf(buf, ":%d", task%step)
				buf.WriteTo(conn)
				buf.Reset()
			}
		}(i)
	}

	timeStart := time.Nanoseconds()
	sentStart := 1
	for sent := 1; sent <= count; sent++ {
		tasks <- sent
		if sent%step == 0 {
			log.Printf("Processed %d packets of %d, QPS=%v", sent, count, float64(sent-sentStart)/float64(time.Nanoseconds()-timeStart)*1e9)
			timeStart = time.Nanoseconds()
			sentStart = sent
		}
	}
	time.Sleep(delay * 2)
}
