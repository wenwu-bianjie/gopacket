package writeJson

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/wenwu-bianjie/gopacket/fasterParsePackets/toKafka/syncProducer"
)

type Vidata struct {
	CaptureFilename string   `json:"capture_filename"`
	Frame           int64    `json:"frame"`
	Time            string   `json:"time"`
	FrameBytes      int      `json:"frame_bytes"`
	SrcMAC          string   `json:"src_mac"`
	DstMAC          string   `json:"dst_mac"`
	SrcIP           string   `json:"src_ip"`
	DstIP           string   `json:"dst_ip"`
	SrcIpv6         string   `json:"src_ipv6"`
	DstIpv6         string   `json:"dst_ipv6"`
	SrcPort         string   `json:"src_port"`
	DstPort         string   `json:"dst_port"`
	IPVersion       string   `json:"ip_version"`
	TCPFlags        TCPFlags `json:"tcp_flags"`
	Identification  uint16   `json:"identification"`
	Seq             uint32   `json:"seq"`
	Ack             uint32   `json:"ack"`
	PayloadBytes    int      `json:"payload_bytes"`
	Payload         string   `json:"payload"`
	SysName         string   `json:"sys_name"`
}

type TCPFlags struct {
	FIN bool `json:"FIN"`
	SYN bool `json:"SYN"`
	RST bool `json:"RST"`
	PSH bool `json:"PSH"`
	ACK bool `json:"ACK"`
	URG bool `json:"URG"`
	ECE bool `json:"ECE"`
	CWR bool `json:"CWR"`
	NS  bool `json:"NS"`
}

var KafkaProducer sarama.SyncProducer
var err error

func init() {
	KafkaProducer, err = syncProducer.NewProducer()
	if err != nil {
		fmt.Println("kafkaProducer ERROR")
		panic(err)
	}
}

func write(tmpdate *Vidata, w *bufio.Writer) error {
	b, err := json.Marshal(&tmpdate)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: "pcap",
		Key:   sarama.StringEncoder(tmpdate.Frame),
	}
	msg.Value = sarama.ByteEncoder(b)

	_, _, err = KafkaProducer.SendMessage(msg)
	fmt.Println(err)
	//_, err = w.Write(b)
	//if err != nil {
	//	fmt.Println("Error", err)
	//	return err
	//}
	//_, err = w.WriteString("\n")
	return err
}
