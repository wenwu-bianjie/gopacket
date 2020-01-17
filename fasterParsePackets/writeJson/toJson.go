package writeJson

import (
	"bufio"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"time"
)

var eth layers.Ethernet
var ip4 layers.IPv4
var ip6 layers.IPv6
var tcp layers.TCP
var payload gopacket.Payload

var ipv4 = "IPv4"
var ipv6 = "IPv6"

var parser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &payload)
var decodedLayers = make([]gopacket.LayerType, 0, 10)

var tcpFlags = TCPFlags{}
var tmp = Vidata{TCPFlags: tcpFlags}

func WritePacketjson(data []byte, w *bufio.Writer, frame int64, fileName string, timestamp time.Time, captureLength int) error {
	tmp.CaptureFilename = fileName
	tmp.Frame = frame
	tmp.Time = timestamp.String()
	tmp.FrameBytes = captureLength
	err := parser.DecodeLayers(data, &decodedLayers)

	for _, typ := range decodedLayers {
		switch typ {
		case layers.LayerTypeEthernet:
			//fmt.Println("    Eth ", eth.SrcMAC, eth.DstMAC)
			tmp.SrcMAC = eth.SrcMAC.String()
			tmp.DstMAC = eth.DstMAC.String()
		case layers.LayerTypeIPv4:
			//fmt.Println("    IP4 ", ip4.SrcIP, ip4.DstIP)
			tmp.SrcIP = ip4.SrcIP.String()
			tmp.DstIP = ip4.DstIP.String()
			tmp.IPVersion = ipv4
			tmp.Identification = ip4.Id
		case layers.LayerTypeIPv6:
			//fmt.Println("    IP6 ", ip6.SrcIP, ip6.DstIP)
			tmp.SrcIpv6 = ip6.SrcIP.String()
			tmp.DstIpv6 = ip6.DstIP.String()
			tmp.IPVersion = ipv6
		case layers.LayerTypeTCP:
			tmp.SrcPort = tcp.SrcPort.String()
			tmp.DstPort = tcp.DstPort.String()
			tmp.Ack = tcp.Ack
			tmp.Seq = tcp.Seq
			tmp.Payload = string(tcp.Payload)
			tmp.PayloadBytes = len(tcp.Payload)
			tmp.TCPFlags.FIN = tcp.FIN
			tmp.TCPFlags.SYN = tcp.SYN
			tmp.TCPFlags.RST = tcp.RST
			tmp.TCPFlags.PSH = tcp.PSH
			tmp.TCPFlags.ACK = tcp.ACK
			tmp.TCPFlags.URG = tcp.URG
			tmp.TCPFlags.ECE = tcp.ECE
			tmp.TCPFlags.CWR = tcp.CWR
			tmp.TCPFlags.NS = tcp.NS
		}
	}

	write(&tmp, w)
	return err
}
