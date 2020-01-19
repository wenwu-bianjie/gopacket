package engine

import (
	"github.com/google/gopacket/pcap"
	"github.com/wenwu-bianjie/gopacket/fasterParsePackets/writeJson"

	"bufio"
	"io"
	"net"
	"strings"
	"syscall"
)

type Scheduler interface {
	Run()
}

type FasterEngine struct {
}

func (e *FasterEngine) Run(handle *pcap.Handle, w *bufio.Writer, fileName string) {
	forEachPackets(handle, w, fileName)
}

func forEachPackets(handle *pcap.Handle, w *bufio.Writer, fileName string) {
	var frame int64

	for {
		data, ci, err := handle.ReadPacketData()
		if err == nil {
			frame++
			writeJson.WritePacketjson(data, w, frame, fileName, ci.Timestamp, ci.CaptureLength)
			continue
		}

		// Immediately retry for temporary network errors
		if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
			continue
		}

		// Immediately retry for EAGAIN
		if err == syscall.EAGAIN {
			continue
		}

		// Immediately break for known unrecoverable errors
		if err == io.EOF || err == io.ErrUnexpectedEOF ||
			err == io.ErrNoProgress || err == io.ErrClosedPipe || err == io.ErrShortBuffer ||
			err == syscall.EBADF ||
			strings.Contains(err.Error(), "use of closed file") {
			break
		}
	}
}
