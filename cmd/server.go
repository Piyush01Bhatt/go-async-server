package main

import (
	"fmt"
	"net"
	"syscall"
)

const (
	address    = "127.0.0.1"
	port       = 8080
	maxClients = 10
)

func main() {
	fmt.Println("Starting Async Server...")

	// Create a socket file descriptor
	sockFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("Error creating socket:", err)
		return
	}
	defer syscall.Close(sockFd)

	// Bind the socket to an address and port
	ip4 := net.ParseIP(address)
	if ip4 == nil {
		fmt.Println("Invalid IPv4 address")
	}
	addr := syscall.SockaddrInet4{Port: port, Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]}}
	if err := syscall.Bind(sockFd, &addr); err != nil {
		fmt.Println("Error binding socket")
		return
	}

	// Listen to incoming connections
	if err := syscall.Listen(sockFd, maxClients); err != nil {
		fmt.Println("Error Listening to socket")
		return
	}

	// Accept connection
	clientFd, _, err := syscall.Accept(sockFd)
	if err != nil {
		fmt.Println("Error accepting client:", err)
		return
	}
	defer syscall.Close(clientFd)

	// Read from connection
	buf := make([]byte, 1024)
	_, err = syscall.Read(clientFd, buf)

	if err != nil {
		fmt.Println("Error Reading from connection")
		return
	}

	fmt.Println(string(buf))
}
