package engine

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Scheduler interface {
	ReadyNotifier
	PacketChan() chan gopacket.Packet
	Run()
	Submit(gopacket.Packet)
}

type ReadyNotifier interface {
	WorkerReady(chan gopacket.Packet)
}

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}

func (e *ConcurrentEngine) Run(handle *pcap.Handle) {
	e.Scheduler.Run()
	out := make(chan *ParseResult)

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorkerQueued(e.Scheduler.PacketChan(), out, e.Scheduler)
	}

	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var packetsLen int

	for packet := range packetSource.Packets() {
		e.Scheduler.Submit(packet)
		packetsLen++
	}

	var count int

	for {
		parseResult := <-out
		count++
		for _, layerType := range parseResult.FoundLayerTypes {
			if layerType == layers.LayerTypeEthernet {
				fmt.Println("Ethernet: ", *parseResult.EthLayer)
			}
			if layerType == layers.LayerTypeIPv4 {
				fmt.Println("IPv4: ", (*parseResult.IpLayer).SrcIP, "->", (*parseResult.IpLayer).DstIP)
			}
			if layerType == layers.LayerTypeIPv6 {
				fmt.Println("IPv6: ", (*parseResult.Ipv6Layer).SrcIP, "->", (*parseResult.Ipv6Layer).DstIP)
			}
			if layerType == layers.LayerTypeTCP {
				fmt.Println("TCP Port: ", (*parseResult.IcpLayer).SrcPort, "->", (*parseResult.IcpLayer).DstPort)
				fmt.Println("TCP SYN:", (*parseResult.IcpLayer).SYN, " | ACK:", (*parseResult.IcpLayer).ACK)
			}
		}

		if count >= packetsLen {
			break
		}
	}
	fmt.Println(count)
}

func (e *ConcurrentEngine) createWorkerQueued(in chan gopacket.Packet, out chan *ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			ready.WorkerReady(in)
			packet := <-in
			result, _ := fasterParsePackets(packet)
			out <- result
		}
	}()
}

func (e *ConcurrentEngine) createWorker(in chan gopacket.Packet, out chan *ParseResult) {
	go func() {
		for {
			packet := <-in
			result, _ := fasterParsePackets(packet)
			out <- result
		}
	}()
}
