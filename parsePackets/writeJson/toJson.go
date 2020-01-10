package writeJson

import (
	"bufio"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
)

var mutex sync.Mutex

var httpMethods = [...]string{"OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "CONNECT"}

func macprocess(mac string) string {
	str := strings.Replace(mac, ":", "", -1)
	return str
}

func findurl(date string) (string, string) {
	urlitems := strings.Split(date, " ")
	if 2 > len(urlitems) {
		fmt.Println("urlitems < 2")
		return "", ""
	}
	var Url string
	var Ua string
	if len(urlitems) > 3 {
		domain := strings.Split(urlitems[3], "\r\n")
		path := strings.Replace(urlitems[1], "\u0026", "&", -1)
		Url = "http://" + domain[0] + path
	}
	if len(urlitems) > 4 {
		tmp := strings.Split(urlitems[4], "\r\n")
		Ua = tmp[0]
	}
	return Url, Ua
}

func printPacketInfo(packet gopacket.Packet) {
	// Let's see if the packet is an ethernet packet
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		fmt.Println("Ethernet layer detected.")
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
		fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
		// Ethernet type is typically IPv4 but could be ARP or other
		fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
		fmt.Println()
	}

	// Let's see if the packet is IP (even though the ether type told us)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		fmt.Println("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)

		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
		fmt.Println("Protocol: ", ip.Protocol)
		fmt.Println()
	}

	// Let's see if the packet is TCP
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		fmt.Println("TCP layer detected.")
		tcp, some := tcpLayer.(*layers.TCP)
		// tcp := tcpLayer
		fmt.Println(reflect.TypeOf(tcp), reflect.TypeOf(tcpLayer))
		// fmt.Println(tcpLayer)
		fmt.Println("some=", some)

		// TCP layer variables:
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		// fmt.Printf("From port %d to %d\n", tcpLayer.SrcPort, tcpLayer.DstPort)
		fmt.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
		fmt.Println("Sequence number: ", tcp.Seq)
		fmt.Println()
	}

	// Iterate over all layers, printing out each layer type
	fmt.Println("All packet layers:")
	for _, layer := range packet.Layers() {
		fmt.Println("- ", layer.LayerType())
	}

	// When iterating through packet.Layers() above,
	// if it lists Payload layer then that is the same as
	// this applicationLayer. applicationLayer contains the payload
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		fmt.Println("Application layer/Payload found.")
		fmt.Printf("%s\n", applicationLayer.Payload())

		// Search for a string inside the payload
		if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
			fmt.Println("HTTP found!")
		}
	}

	// Check for errors
	if err := packet.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}
}

func WritePacketjson(packet gopacket.Packet, file_fd *os.File) error {
	tmpdate := Vidata{}
	// Let's see if the packet is an ethernet packet
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		tmpdate.SrcMAC = ethernetPacket.SrcMAC.String()
		tmpdate.DstMAC = ethernetPacket.DstMAC.String()
	}

	// Let's see if the packet is IP (even though the ether type told us)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		tmpdate.SrcIP = ip.SrcIP.String()
		tmpdate.DstIP = ip.DstIP.String()
	}

	ipLayer = packet.Layer(layers.LayerTypeIPv6)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv6)
		tmpdate.SrcIP = ip.SrcIP.String()
		tmpdate.DstIP = ip.DstIP.String()
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		tmpdate.SrcPort = tcp.SrcPort.String()
		tmpdate.DstPort = tcp.DstPort.String()
	}

	// When iterating through packet.Layers() above,
	// if it lists Payload layer then that is the same as
	// this applicationLayer. applicationLayer contains the payload
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		//fmt.Println("Application layer/Payload found.")

		// Search for a string inside the payload
		payloadStr := string(applicationLayer.Payload())
		for _, verb := range httpMethods {
			if strings.Contains(payloadStr, verb) {
				tmpdate.Url, tmpdate.Ua = findurl(string(applicationLayer.Payload()))
				tmpdate.HttpRequest = strings.ReplaceAll(string(applicationLayer.Payload()), "\r\n", "  ")
				break
			}
		}

		if strings.Contains(payloadStr, "Content-Type") {
			tmpdate.Url, tmpdate.Ua = findurl(string(applicationLayer.Payload()))
			response := strings.Split(string(applicationLayer.Payload()), "\r\n\r\n")
			if len(response) > 0 {
				fmt.Printf("%v\n", response)
				tmpdate.HttpResponse = strings.ReplaceAll(response[0], "\r\n", "  ")
			}
		}
	}

	// Check for errors
	if err := packet.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	err := write(tmpdate, file_fd)
	return err
}

func determineEncoding(
	r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(
		bytes, "")
	return e
}
