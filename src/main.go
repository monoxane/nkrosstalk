package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	// Mess of env vars and converstions to ints
	listenPort := getEnv("NK_LISTEN", "7788")
	nkHost := getEnv("NK_HOST", "10.0.0.0")
	nkAddress, _ := strconv.ParseInt(getEnv("NK_SIZE", "254"), 10, 0)
	nkSize, _ := strconv.ParseInt(getEnv("NK_SIZE", "16"), 10, 0)

	// Suprise suprise more type casting
	Router := IPS{
		Host:         nkHost,
		Address:      uint8(nkAddress),
		Destinations: uint16(nkSize),
		Sources:      uint16(nkSize),
	}

	go Router.Connect() // Fork it out here or it hangs

	listener, err := net.Listen("tcp", "0.0.0.0:"+listenPort)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	// Let the user know what's going on
	log.Println("Listening for RossTalk on tcp://0.0.0.0:" + listenPort)
	log.Println("XPT <level>:<destination>:<source>")

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleClientRequest(con, Router) // Fork request process
	}
}

func handleClientRequest(con net.Conn, r IPS) {
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
