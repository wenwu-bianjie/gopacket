package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"flag"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

var (
	// deviceName  string = "eth0"
	snapshotLen int32         = 1024
	promiscuous bool          = false
	timeout     time.Duration = -1 * time.Second
)

var (
	pcapFile string = "./test4.pcap"
)

var filter = flag.String("f", "tcp", "BPF filter for pcap")

func main() {
	// Find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// Print device information
	fmt.Println("Devices found:")
	f, _ := os.Create(pcapFile)
	w := pcapgo.NewWriter(f)
	w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet)
	defer f.Close()

	packetChan := make(chan gopacket.Packet)
	last := len(devices) - 1
	d := devices[last]
	go getPacket(d.Name, packetChan)

	//for _, d := range devices {
	//	fmt.Println("\nName: ", d.Name)
	//	fmt.Println("Description: ", d.Description)
	//	fmt.Println("Devices addresses: ", d.Description)
	//	for _, address := range d.Addresses {
	//		fmt.Println("- IP address: ", address.IP)
	//		fmt.Println("- Subnet mask: ", address.Netmask)
	//	}
	//}

	for {
		select {
		case packet := <-packetChan:
			// Process packet here
			w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		}
	}
}

func getPacket(deviceName string, packetChan chan gopacket.Packet) {
	// Open output pcap file and write header

	// Open the device for capturing
	handle, err := pcap.OpenLive(deviceName, snapshotLen, promiscuous, timeout)
	if err := handle.SetBPFFilter(*filter); err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	if err != nil {
		fmt.Printf("Error opening device %s: %v", deviceName, err)
	} else {
		// Start processing packets
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for {
			packetSourceChane := packetSource.Packets()
			select {
			case packet := <-packetSourceChane:
				packetChan <- packet
			}
		}
	}
}
