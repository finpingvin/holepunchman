package main

import "net"

type Peer struct {
	Addr *net.UDPAddr
	Conn *net.UDPConn
}

type PeerInfo struct {
	IP   string
	Port string
}
