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
		if command[:len(command)-1] == "exit 0" {
			os.Exit(0)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		fmt.Println(command[:len(command)-1] + ": command not found")
	}

	// Wait for user input
	//bufio.NewReader(os.Stdin).ReadString('\n')
}
