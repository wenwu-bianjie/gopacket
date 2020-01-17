package scheduler

import "github.com/google/gopacket"

type SimpleScheduler struct {
	packetChan chan gopacket.Packet
}

func (s *SimpleScheduler) PacketChan() chan gopacket.Packet {
	return s.packetChan
}

func (s *SimpleScheduler) Run() {
	s.packetChan = make(chan gopacket.Packet)
}

func (s *SimpleScheduler) Submit(p gopacket.Packet) {
	go func() {
		s.packetChan <- p
	}()
}
