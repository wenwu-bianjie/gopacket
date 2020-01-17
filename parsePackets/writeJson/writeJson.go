package writeJson

import (
	"encoding/json"
	"fmt"
	"os"
)

type Vidata struct {
	CaptureFilename string `json:"capture_filename"`
	Frame           int64  `json:"frame"`
	Time            string `json:"time"`
	FrameBytes      int64  `json:"frame_bytes"`
	SrcMAC          string `json:"src_mac"`
	DstMAC          string `json:"dst_mac"`
	SrcIP           string `json:"src_ip"`
	DstIP           string `json:"dst_ip"`
	SrcIpv6         string `json:"src_ipv6"`
	DstIpv6         string `json:"dst_ipv6"`
	SrcPort         string `json:"src_port"`
	DstPort         string `json:"dst_port"`
	IPVersion       string `json:"ip_version"`
	TCPFlags        string `json:"tcp_flags"`
	Identification  int    `json:"identification"`
	Seq             uint32 `json:"seq"`
	Ack             uint32 `json:"ack"`
	PayloadBytes    int    `json:"payload_bytes"`
	Payload         string `json:"payload"`
	SysName         string `json:"sys_name"`
}

func write(tmpdate Vidata, file_fd *os.File) error {
	b, err := json.Marshal(tmpdate)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	_, err = file_fd.Write(b)
	if err != nil {
		fmt.Println("Error", err)
		return err
	}
	_, err = file_fd.WriteString("\n")
	return err
}
