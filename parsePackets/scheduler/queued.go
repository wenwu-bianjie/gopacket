package scheduler

import "github.com/google/gopacket"

type QueuedScheduler struct {
	packetChan chan gopacket.Packet
	workerChan chan chan gopacket.Packet
}

func (s *QueuedScheduler) PacketChan() chan gopacket.Packet {
	return make(chan gopacket.Packet)
}

func (s *QueuedScheduler) Submit(r gopacket.Packet) {
	s.packetChan <- r
}

func (s *QueuedScheduler) WorkerReady(
	w chan gopacket.Packet) {
	s.workerChan <- w
}

func (s *QueuedScheduler) Run() {
	s.workerChan = make(chan chan gopacket.Packet)
	s.packetChan = make(chan gopacket.Packet)
	go func() {
		var packetQ []gopacket.Packet
		var workerQ []chan gopacket.Packet

		for {
			var activePacket gopacket.Packet
			var activeWorker chan gopacket.Packet
			if len(packetQ) > 0 && len(workerQ) > 0 {
				activePacket = packetQ[0]
				activeWorker = workerQ[0]
			}

			select {
			case p := <-s.packetChan:
				packetQ = append(packetQ, p)
			case w := <-s.workerChan:
				workerQ = append(workerQ, w)
			case activeWorker <- activePacket:
				packetQ = packetQ[1:]
				workerQ = workerQ[1:]
			}
		}
	}()
}
