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
var builtin = []string{"exit", "echo", "type"}

func process_command(command string) (string, string) {
	if strings.Count(command, " ") >= 1 {
		chunks := strings.SplitN(command, " ", 2)
		return chunks[0], chunks[1][:len(chunks[1])-1]
	}
	return command[:len(command)-1], ""
}

func main() {
	for i := 0; i < 100; i++ {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		cmd, arg := process_command(command)

		switch cmd {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(arg)
		case "type":
			if slices.Contains(builtin, arg) {
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
