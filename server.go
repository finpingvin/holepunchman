package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func RunServer(serverAddr *string) {
	addr := fmt.Sprintf(":%s", *serverAddr)
	udpAddr, _ := net.ResolveUDPAddr("udp4", addr)
	conn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Relay server listening on", addr)

	var peer1, peer2 *net.UDPAddr

	buf := make([]byte, 1024)
	for {
		n, remoteAddr, _ := conn.ReadFromUDP(buf)
		message := string(buf[:n])
		fmt.Println("Received from", remoteAddr, ":", message)

		if message == "register" {
			if peer1 == nil {
				peer1 = remoteAddr
				continue
			}
			if peer2 == nil {
				peer2 = remoteAddr
			}

			// Send peer info
			info1 := PeerInfo{IP: peer1.IP.String(), Port: fmt.Sprintf("%d", peer1.Port)}
			info2 := PeerInfo{IP: peer2.IP.String(), Port: fmt.Sprintf("%d", peer2.Port)}
			data1, _ := json.Marshal(info2)
			data2, _ := json.Marshal(info1)

			conn.WriteToUDP(data1, peer1)
			conn.WriteToUDP(data2, peer2)

			// Send "start punching" message to both
			conn.WriteToUDP([]byte("START"), peer1)
			conn.WriteToUDP([]byte("START"), peer2)

			fmt.Println("Sent peer info + START to both clients")

			// Reset for next pair
			peer1, peer2 = nil, nil
		}
	}
}
