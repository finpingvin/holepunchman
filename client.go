package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func RunClient(serverIp *string, serverPort *string) {
	serverAddr := fmt.Sprintf("%s:%s", *serverIp, *serverPort)
	localConn, err := net.ListenUDP("udp4", nil)
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

	// Jab a bit
	for i := 0; i < 30; i++ {
		msg := fmt.Sprintf("punch %d", i)
		localConn.WriteToUDP([]byte(msg), peerUDPAddr)
		time.Sleep(200 * time.Millisecond)
	}

	// Punch receiver goroutine
	go func() {
		fmt.Println("Reading messages")
		for {
			n, addr, err := localConn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("Read error:", err)
				continue
			}
			if n == 0 {
				continue
			}
			msg := string(buf[:n])
			fmt.Printf("Received from %s: %s\n", addr, msg)
		}
	}()

	// Puncher goroutine
	go func() {
		fmt.Printf("Writing messages to: %s:%s\n", peer.IP, peer.Port)
		for {
			msg := fmt.Sprintf("Punch to: %s", peer.Port)
			localConn.WriteToUDP([]byte(msg), peerUDPAddr)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// Trade punches for 60 seconds
	time.Sleep(60 * time.Second)
	fmt.Println("Done. Exiting.")
}
