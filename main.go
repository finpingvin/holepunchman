package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	mode := flag.String("mode", "", "Mode to run: 'server' or 'client'")
	serverIp := flag.String("ip", "", "Server address (client mode only)")
	serverPort := flag.String("port", "", "Server port (client mode only)")
	flag.Parse()

	switch *mode {
	case "server":
		fmt.Println("Running in server mode")
		RunServer(serverPort)

	case "client":
		if *serverIp == "" || *serverPort == "" {
			fmt.Println("Client mode requires -ip and -port flags")
			os.Exit(1)
		}
		fmt.Printf("Running in client mode, connecting to %s:%s\n", *serverIp, *serverPort)
		RunClient(serverIp, serverPort)

	default:
		fmt.Println("Invalid mode. Use -mode=server or -mode=client")
		flag.Usage()
		os.Exit(1)
	}
}
