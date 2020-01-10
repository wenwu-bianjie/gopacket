package writeJson

import (
	"encoding/json"
	"fmt"
	"os"
)

type Vidata struct {
	//State0       string `json:"state0, string"`
	//State1       string `json:"state1, string"`
	SrcMAC       string `json:"src_mac"`
	DstMAC       string `json:"dst_mac"`
	SrcIP        string `json:"src_ip"`
	DstIP        string `json:"dst_ip"`
	SrcIpv6      string `json:"src_ipv6"`
	DstIpv6      string `json:"dst_ipv6"`
	SrcPort      string `json:"src_port"`
	DstPort      string `json:"dst_port"`
	Url          string `json:"url, string"`
	Ua           string `json:"ua, string"`
	HttpRequest  string `json:"http_request"`
	HttpResponse string `json:"http_response"`
	//Reffer string `json:"reffer, string"`
	//Cookie string `json:"cookie, string"`
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
