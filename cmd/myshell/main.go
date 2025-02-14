package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
// var _ = fmt.Fprint
var builtin = []string{"exit", "echo"}

func main() {
	// Uncomment this block to pass the first stage

	for i := 0; i < 100; i++ {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		chunks := strings.Split(command, " ")
		cmd, arg := chunks[0], chunks[1]
		switch cmd {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(command[5 : len(command)-1])
		case "type":
			if slices.Contains(builtin, command[5:len(command)-1]) {
				fmt.Println(arg + " is a shell builtin")
			} else {
				fmt.Println(arg + ": not found")
			}
		default:
			fmt.Println(command[:len(command)-1] + ": command not found")
		}

		if err != nil {
			//panic(err)
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

	}

	// Wait for user input
	//bufio.NewReader(os.Stdin).ReadString('\n')
}
