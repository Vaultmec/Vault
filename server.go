package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	clients     = make(map[net.Conn]string)
	clientsLock sync.Mutex
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on :8080")

	// Handle graceful shutdown on Ctrl+C
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		fmt.Println("\nServer shutting down...")
		listener.Close()

		// Notify all clients to disconnect
		clientsLock.Lock()
		for conn := range clients {
			conn.Close()
		}
		clientsLock.Unlock()

		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Prompt the client for their name

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	clientName := scanner.Text()

	clientsLock.Lock()
	clients[conn] = clientName
	clientsLock.Unlock()

	// Welcome message to the client
	fmt.Fprintf(conn, "Welcome, %s!\n", clientName)
	fmt.Printf("Client %s connected.\n", clientName)

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Client %s disconnected.\n", clientName)
			clientsLock.Lock()
			delete(clients, conn)
			clientsLock.Unlock()
			return
		}

		message := string(buffer[:n])
		fmt.Printf("[%s]: %s", clientName, message)

		// Broadcast the message to all clients
		clientsLock.Lock()
		for otherConn := range clients {
			if otherConn != conn {
				otherConn.Write([]byte(fmt.Sprintf("[%s]: %s", clientName, message)))
			}
		}
		clientsLock.Unlock()
	}
}
