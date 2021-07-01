package nk

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type NKType struct {
	Host         string
	Address      uint8
	Destinations uint16
	Sources      uint16
	Conn         net.Conn
}

func crc16(buffer []byte) uint16 {
	var crc = 0xFFFF
	var odd = 0x0000

	for i := 0; i < len(buffer); i++ {
		crc = crc ^ int(buffer[i])

		for j := 0; j < 8; j++ {
			odd = crc & 0x0001
			crc = crc >> 1
			if odd == 1 {
				crc = crc ^ 0xA001
			}
		}
	}

	crc = ((crc & 0xFF) << 8) | ((crc & 0xFF00) >> 8)
	return uint16(crc)
}

// XPT Just returns payload to send to router to close xpt
func (n *NKType) GenerateXPTRequest(level uint32, destination uint16, source uint16) ([]byte, error) {
	if !(level < 8) && !(destination <= n.Destinations) && !(source <= n.Sources) {
		return []byte{}, errors.New("source or destination out of bounds")
	} else {
		// log.Println("Data:", level, destination, source, n.Address)

		var destination uint16 = destination - 1
		var source uint16 = source - 1

		type TBusPacketPayload struct {
			NK2Header   uint32
			RTRAddress  uint8
			UNKNB       uint16
			Destination uint16
			Source      uint16
			LevelMask   uint32
			UNKNC       uint8
		}

		type TBusPacket struct {
			HeaderA uint32
			HeaderB uint16
			Payload TBusPacketPayload
			CRC     uint16
		}

		payload := TBusPacketPayload{
			NK2Header:   0x4e4b3200,
			RTRAddress:  n.Address,
			UNKNB:       0x0409,
			Destination: destination,
			Source:      source,
			LevelMask:   level,
			UNKNC:       0x00,
		}

		payloadBuffer := new(bytes.Buffer)
		err := binary.Write(payloadBuffer, binary.BigEndian, payload)
		if err != nil {
			log.Println("TBusPacketPayload binary.Write failed:", err)
		}

		packet := TBusPacket{
			HeaderA: 0x50415332,
			HeaderB: 0x0012,
			Payload: payload,
			CRC:     crc16(payloadBuffer.Bytes()),
		}

		// log.Println(packet)
		packetBuffer := new(bytes.Buffer)
		err = binary.Write(packetBuffer, binary.BigEndian, packet)
		if err != nil {
			log.Println("TBustPacket binary.Write failed:", err)
		}
		// log.Printf("%x", packetBuffer.Bytes())
		// log.Println("5041533200124e4b3200fe04090040003100000001005de9")

		return packetBuffer.Bytes(), nil
	}
}

func (n *NKType) Connect() {
	conn, err := net.Dial("tcp", n.Host+":5000")
	if err != nil {
		log.Fatalln(err)
	}
	n.Conn = conn
	defer n.Conn.Close()

	serverReader := bufio.NewReader(n.Conn)

	openConnStr := strings.TrimSpace("PHEONIX-DB")
	if _, err = n.Conn.Write([]byte(openConnStr + "\n")); err != nil {
		log.Printf("failed to send the client request: %v\n", err)
	}

	go func() {
		for range time.Tick(10 * time.Second) {
			n.Conn.Write([]byte("HI"))
		}
	}()

	for {

		serverResponse, err := serverReader.ReadString('\n')
		switch err {
		case nil:
			log.Println(strings.TrimSpace(serverResponse))
		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}
	}
}

func (n *NKType) SetCrosspoint(level uint32, destination uint16, source uint16) error {
	XPTRequest, err := n.GenerateXPTRequest(level, destination, source)
	if err != nil {
		return err
	}
	_, err = n.Conn.Write(XPTRequest)
	if err != nil {
		return err
	}
	return nil
}

// 5041533200124e4b3200fe04090040003100000001005de9
// 504153320012

// h := (uint32)(((fh.year*100+fh.month)*100+fh.day)*100 + fh.h)
// a := make([]byte, 4)
// binary.LittleEndian.PutUint32(a, h)

// Do it as a struct

// type Test struct {
// 	Byte1 byte
// 	Byte2 byte
// 	ByteArr1 [2]byte
// 	My16 uint16
// }

// t := Test{}
// err = binary.Read(bytes.NewReader(buff), binary.BigEndian, &t)
// if err != nil {
// 	fmt.Println("binary.Read failed:", err)
// }
// //buff = []byte{}
// log.Print(t)

// var pi float64
// b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
// buf := bytes.NewReader(b)
// err := binary.Read(buf, binary.LittleEndian, &pi)
// if err != nil {
// 		fmt.Println("binary.Read failed:", err)
// }
// fmt.Print(pi)
