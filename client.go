package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type PeerInfo struct {
	IP   string
	Port string
}

func main() {
	serverAddr := "my.server.ip:9000" // Replace with your server IP
	localConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(err)
	}
	defer localConn.Close()

	localAddr := localConn.LocalAddr().(*net.UDPAddr)
	fmt.Println("Local UDP address:", localAddr)

	// Register with server
	serverUDPAddr, _ := net.ResolveUDPAddr("udp", serverAddr)
	_, err = localConn.WriteToUDP([]byte("register"), serverUDPAddr)
	if err != nil {
		panic(err)
	}

	// Wait for peer info
	buf := make([]byte, 1024)
	n, _, err := localConn.ReadFromUDP(buf)
	if err != nil {
		panic(err)
	}
	var peer PeerInfo
	if err := json.Unmarshal(buf[:n], &peer); err != nil {
		panic(err)
	}
	fmt.Println("Received peer info:", peer)

	peerAddrStr := net.JoinHostPort(peer.IP, peer.Port)
	peerUDPAddr, _ := net.ResolveUDPAddr("udp", peerAddrStr)

	// Punch hole
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("punch %d", i)
		localConn.WriteToUDP([]byte(msg), peerUDPAddr)
		time.Sleep(200 * time.Millisecond)
	}

	// Decide role: initiator or responder
	peerPortInt, err := strconv.Atoi(peer.Port)
	if err != nil {
		panic(err)
	}
	isInitiator := min(localAddr.Port, peerPortInt) == localAddr.Port

	fmt.Println("Initiator?", isInitiator, localAddr.Port, peerPortInt)

	// Listener goroutine
	go func() {
		for {
			fmt.Println("Reading messages")
			n, addr, err := localConn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("Read error:", err)
				continue
			}
			msg := string(buf[:n])
			fmt.Printf("Received from %s: %s\n", addr, msg)

			if strings.TrimSpace(msg) == "ping" {
				localConn.WriteToUDP([]byte("pong"), peerUDPAddr)
				fmt.Println("Sent pong")
			}
		}
	}()

	if isInitiator {
		// Initiator sends ping
		time.Sleep(1 * time.Second)
		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("ping %d", i)
			fmt.Println("Sending ping to peer...", peerUDPAddr)
			localConn.WriteToUDP([]byte(msg), peerUDPAddr)
			time.Sleep(200 * time.Millisecond)
		}
		// localConn.WriteToUDP([]byte("ping"), peerUDPAddr)
	}

	// Wait to allow ping/pong exchange to complete
	time.Sleep(30 * time.Second)
	fmt.Println("Done. Exiting.")
}
