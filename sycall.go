package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Choose a system call:")
	fmt.Println("1. ls (list directory contents)")
	fmt.Println("2. pwd (print working directory)")
	fmt.Println("3. date (display the current date and time)")
	fmt.Println("4. uname (print system information)")

	var choice int
	fmt.Print("Enter your choice (1-4): ")
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch choice {
	case 1:
		runCommand("ls")
	case 2:
		runCommand("pwd")
	case 3:
		runCommand("date")
	case 4:
		runCommand("uname", "-a")
	default:
		fmt.Println("Invalid choice. Please enter a number between 1 and 4.")
	}
}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
