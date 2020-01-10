package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"github.com/wenwu-bianjie/gopacket/parsePackets/engine"
	"github.com/wenwu-bianjie/gopacket/parsePackets/scheduler"

	"log"
	"os"
	"time"
)

var (
	pcapFile string = "./test.pcap"
	handle   *pcap.Handle
	err      error
	filename string = "./test.json"
)

func main() {
	t1 := time.Now()
	// Open file instead of device
	handle, err = pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	file_fd, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Cannot open file %s!\n", filename)
		return
	}
	defer file_fd.Close()

	e := engine.ToJsonEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 1,
	}
	e.Run(handle, file_fd)
	t2 := time.Now()
	fmt.Println(t2.Sub(t1))
}
