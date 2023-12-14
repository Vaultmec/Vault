
import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// Create channels for communication
	messageChannel := make(chan string)
	doneChannel := make(chan struct{})

	// Start the sender and receiver goroutines
	wg.Add(2)
	go sender("User 1", messageChannel, doneChannel, &wg)
	go receiver("User 2", messageChannel, doneChannel, &wg)

	// Wait for both goroutines to finish
	wg.Wait()
}

func sender(name string, messageChannel chan<- string, doneChannel <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		var message string

		// Get input from the sender
		fmt.Printf("[%s] Enter a message (type 'exit' to quit): ", name)
		fmt.Scanln(&message)

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
