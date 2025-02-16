package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
// var _ = fmt.Fprint
var builtin = []string{"exit", "echo", "type"}
var PATH = os.Getenv("PATH")

func process_command(command string) (string, string) {
	if strings.Count(command, " ") >= 1 {
		chunks := strings.SplitN(command, " ", 2)
		return chunks[0], strings.TrimSpace(chunks[1][:len(chunks[1])-1])
	}
	return command[:len(command)-1], ""
}

func main() {
	for i := 0; i < 100; i++ {
		fmt.Fprint(os.Stdout, "$ ")
		full_command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		command_keyword, arg := process_command(full_command)
		switch command_keyword {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(arg)
		case "type":
			if slices.Contains(builtin, arg) {
				fmt.Println(arg + " is a shell builtin")
			} else {
				dirs := strings.Split(PATH, string(filepath.ListSeparator))
				for i := 0; i < len(dirs); i++ {
					search_path := dirs[i] + "\\" + arg
					if _, search_err := os.Stat(search_path); search_err == nil {
						fmt.Println("valid_command is " + search_path)
						break
					} else if i+1 == len(dirs) {
						fmt.Println(arg + ": not found")
					}

				}

			}

		default:
			fmt.Println(full_command[:len(full_command)-1] + ": command not found")
		}

		if err != nil {
			//panic(err)
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

	}
}
