package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage

	for i := 0; i < 100; i++ {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		switch cmd := command[:4]; cmd {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(command[5 : len(command)-1])
		default:
			fmt.Println(command[:len(command)-1] + ": command not found")
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

	}

	// Wait for user input
	//bufio.NewReader(os.Stdin).ReadString('\n')
}
