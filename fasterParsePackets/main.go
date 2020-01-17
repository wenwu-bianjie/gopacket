package main

import (
	"bufio"
	"fmt"
	"github.com/google/gopacket/pcap"
	"github.com/wenwu-bianjie/gopacket/fasterParsePackets/engine"
	"github.com/wenwu-bianjie/gopacket/fasterParsePackets/writeJson"
	"log"
	"os"
	"time"
)

var (
	pcapFile string = "./test3.pcap"
	handle   *pcap.Handle
	err      error
	filename string = "./test3.json"
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

	newWriter := bufio.NewWriterSize(file_fd, 1024*1024)
	defer newWriter.Flush()
	defer writeJson.KafkaProducer.Close()

	e := engine.FasterEngine{}
	e.Run(handle, newWriter, file_fd.Name())

	t2 := time.Now()
	fmt.Println(t2.Sub(t1))
}
