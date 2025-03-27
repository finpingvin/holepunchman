package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type PeerInfo struct {
	IP   string
	Port string
}

var (
	peers []PeerInfo
	mu    sync.Mutex
)

func main() {
	addr := ":9000"
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("NAT relay server listening on", addr)

	buf := make([]byte, 1024)
	for {
		n, remoteAddr, _ := conn.ReadFromUDP(buf)
		message := string(buf[:n])
		fmt.Println("Received:", message, "from", remoteAddr)

		mu.Lock()
		peer := PeerInfo{IP: remoteAddr.IP.String(), Port: fmt.Sprintf("%d", remoteAddr.Port)}
		peers = append(peers, peer)
		fmt.Println("Peers", peers)

		if len(peers) == 2 {
			// Send each other's address to both clients
			p1 := peers[0]
			p2 := peers[1]

			msg1, _ := json.Marshal(p2)
			msg2, _ := json.Marshal(p1)

			addr1, _ := net.ResolveUDPAddr("udp", net.JoinHostPort(p1.IP, p1.Port))
			addr2, _ := net.ResolveUDPAddr("udp", net.JoinHostPort(p2.IP, p2.Port))

			conn.WriteToUDP(msg1, addr1)
			conn.WriteToUDP(msg2, addr2)

			fmt.Println("Exchanged peer info")
			peers = nil // Reset for new pair
		}
		mu.Unlock()
	}
}
