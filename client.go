package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine to receive messages from the server
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Goroutine to send messages to the server
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			message := scanner.Text()
			conn.Write([]byte(message + "\n"))
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()
}
