package engine

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/wenwu-bianjie/gopacket/parsePackets/writeJson"
	"os"
)

type ToJsonEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}

func (e *ToJsonEngine) Run(handle *pcap.Handle, file_fd *os.File) {
	e.Scheduler.Run()
	//out := make(chan error)
	//
	//for i := 0; i < e.WorkerCount; i++ {
	//	e.createWorkerQueued(e.Scheduler.PacketChan(), out, e.Scheduler, file_fd)
	//}

	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var packetsLen int

	for packet := range packetSource.Packets() {
		//e.Scheduler.Submit(packet)
		writeJson.WritePacketjson(packet, file_fd)
		packetsLen++
	}

	fmt.Printf("总数 %v\n", packetsLen)

	//var count int
	//
	//for {
	//	err := <-out
	//	count++
	//
	//	if err != nil {
	//		fmt.Println("Error:", err)
	//	}
	//
	//	if count >= packetsLen {
	//		break
	//	}
	//}
	//fmt.Println(count)
}
