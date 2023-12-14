package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// Create channels for communication
	doneChannel := make(chan struct{})

	// Get names for User 1 and User 2
	fmt.Print("Enter name for User 1: ")
	user1Name := getUserInput()

	fmt.Print("Enter name for User 2: ")
	user2Name := getUserInput()

	// Create channels and start goroutines for User 1
	messageChannelUser1 := make(chan string)
	wg.Add(2)
	go sender(user1Name, messageChannelUser1, doneChannel, &wg)
	go receiver(user1Name, messageChannelUser1, doneChannel, &wg)

	// Create channels and start goroutines for User 2
	messageChannelUser2 := make(chan string)
	wg.Add(2)
	go sender(user2Name, messageChannelUser2, doneChannel, &wg)
	go receiver(user2Name, messageChannelUser2, doneChannel, &wg)

	// Wait for all goroutines to finish
	wg.Wait()
}

func sender(name string, messageChannel chan<- string, doneChannel chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		var message string

		// Get input from the sender
		fmt.Printf("[%s] Enter a message (type 'exit' to quit): ", name)
		message = getUserInput()

		// Send the message to the receiver
		messageChannel <- fmt.Sprintf("[%s]: %s", name, message)

		// Check if the sender wants to exit
		if message == "exit" {
			close(doneChannel)
			return
		}
	}
}

func receiver(name string, messageChannel <-chan string, doneChannel <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case message := <-messageChannel:
			// Receive and display messages from the sender
			fmt.Println(message)

		case <-doneChannel:
			// The sender has exited; exit the receiver
			fmt.Printf("[%s] Exiting...\n", name)
			return
		}
	}
}

func getUserInput() string {
	var input string
	fmt.Scanln(&input)
	return input
}
