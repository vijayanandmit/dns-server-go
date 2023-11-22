package main

import (
	"encoding/binary"
	"fmt"

	// Uncomment this block to pass the first stage

	"net"
)

type Message struct { //12 bytes long
	id     uint16
	qr     bool
	opcode uint8
	aa     bool
	tc     bool   //truncation
	rd     bool   //recursion desired
	ra     bool   //recursion available
	z      uint8  // 3bits reserved, used by DNS SEC queries
	rc     uint8  // 4bits response code
	qc     uint16 //question count
	anc    uint16 //answer record count
	nsc    uint16 //authority record count
	arc    uint16 //additional record count
}

func (m *Message) bytes() []byte {
	out := []byte{}

	out = binary.BigEndian.AppendUint16(out, m.id)

	var flags uint16
	if m.qr {
		flags |= uint16(1) << 15
	}
	flags |= uint16(m.opcode) << 11
	if m.aa {
		flags |= 1 << 10
	}
	if m.tc {
		flags |= 1 << 9
	}
	if m.rd {
		flags |= 1 << 8
	}
	if m.ra {
		flags |= 1 << 7
	}
	flags |= uint16(m.z) << 4
	flags |= uint16(m.rc)
	out = binary.BigEndian.AppendUint16(out, flags)
	out = binary.BigEndian.AppendUint16(out, m.qc)
	out = binary.BigEndian.AppendUint16(out, m.anc)
	out = binary.BigEndian.AppendUint16(out, m.nsc)
	out = binary.BigEndian.AppendUint16(out, m.arc)

	return out
}
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	//	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		if size < 8 {
			fmt.Print("received packet is too small for valid UDP header")
			continue
		}

		srcPort := binary.BigEndian.Uint16(buf[0:2])
		dstPort := binary.BigEndian.Uint16(buf[2:4])
		length := binary.BigEndian.Uint16(buf[4:6])
		checkSum := binary.BigEndian.Uint16(buf[6:8])

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		header := Message{
			id:     1234,
			qr:     false,
			opcode: 0,
			aa:     false,
			tc:     false,
			rd:     false,
			ra:     false,
			z:      0,
			rc:     0,
			qc:     0,
			anc:    0,
			nsc:    0,
			arc:    0,
		}

		response := header.bytes()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
