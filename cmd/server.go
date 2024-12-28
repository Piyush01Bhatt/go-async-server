//go:build linux
// +build linux

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

	// Set socket operation to non-blocking mode
	// if err := syscall.SetNonblock(sockFd, true); err != nil {
	// 	fmt.Println("Error setting socket to non-blocking mode", err)
	// 	return
	// }

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

	// Async IO

	// Create a Epoll File Descriptor
	epollFd, err := syscall.EpollCreate1(0)
	if err != nil {
		fmt.Println("Error creating epoll instance:", err)
		return
	}
	defer syscall.Close(epollFd)

	// Register the server socket with epoll for read events
	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN, // Read event
		Fd:     int32(sockFd),
	}
	if err := syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, sockFd, &event); err != nil {
		fmt.Println("Error adding socket to epoll:", err)
		return
	}

	// We can handle up to 10 events
	events := make([]syscall.EpollEvent, 10)

	fmt.Println("Starting Event Loop...")
	// Event loop
	for {
		// Wait for events
		n, err := syscall.EpollWait(epollFd, events, -1) // -1 means waiting indefinitely
		if err != nil {
			fmt.Println("Error waiting for epoll events:", err)
			return
		}

		// Process the events
		for i := 0; i < n; i++ {
			ev := events[i]

			if ev.Fd == int32(sockFd) {
				// New client connection available
				clientFd, _, err := syscall.Accept(sockFd)
				if err != nil {
					fmt.Println("Error accepting client:", err)
					continue
				}
				fmt.Println("Accepted new client connection.")

				// Register the new client socket with epoll
				clientEvent := syscall.EpollEvent{
					Events: syscall.EPOLLIN, // Read event
					Fd:     int32(clientFd),
				}
				if err := syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, clientFd, &clientEvent); err != nil {
					fmt.Println("Error adding client socket to epoll:", err)
					syscall.Close(clientFd)
					continue
				}
			} else {
				clientFd := int(ev.Fd)
				buf := make([]byte, 1024)

				n, err := syscall.Read(clientFd, buf)
				if err != nil {
					fmt.Println("Error reading from client:", err)
					syscall.Close(clientFd)
					continue
				}

				if n == 0 {
					// Client disconnected
					fmt.Println("Client disconnected.")
					syscall.Close(clientFd)
					continue
				}

				// Echo the received data back to the client
				fmt.Printf("Received data from client: %s\n", string(buf[:n]))
				_, err = syscall.Write(clientFd, buf[:n])
				if err != nil {
					fmt.Println("Error writing to client:", err)
					syscall.Close(clientFd)
					continue
				}

				// Closing the connection
				syscall.Close(clientFd)
			}
		}
	}

}
