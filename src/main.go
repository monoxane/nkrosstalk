package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/monoxane/nkrosstalk/nk"
)

func main() {
	//Router := nk.New("10.101.41.2", 254, 72, 72)
	Router := nk.NKType{
		Host:         "10.101.41.2",
		Address:      254,
		Destinations: 72,
		Sources:      72,
	}

	go Router.Connect()

	listener, err := net.Listen("tcp", "0.0.0.0:9999") // Listen on port 9999 for dev reasons
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleClientRequest(con, Router) // Fork request process
	}
}

func handleClientRequest(con net.Conn, r nk.NKType) {
	defer con.Close()

	clientReader := bufio.NewReader(con)

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')

		if err == nil {
			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			} else {
				command := strings.Split(clientRequest, " ")
				switch command[0] {
				case "XPT": // Crosspoint Request
					level, _ := strconv.ParseInt(strings.Split(command[1], ":")[0], 10, 0)
					destination, _ := strconv.ParseInt(strings.Split(command[1], ":")[1], 10, 0)
					source, _ := strconv.ParseInt(strings.Split(command[1], ":")[2], 10, 8)
					log.Println("Routing", source, "=>", destination, "@", level)
					err := r.SetCrosspoint(uint32(level), uint16(destination), uint16(source))
					if err != nil {
						log.Println(err)
					}
					return
				}
			}
		} else if err == io.EOF {
			log.Println("client closed the connection by terminating the process")
			return
		} else {
			log.Printf("error: %v\n", err)
			return
		}
	}
}
