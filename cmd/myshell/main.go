package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
// var _ = fmt.Fprint
var builtin = []string{"exit", "echo", "type", "pwd", "cd"}
var PATH = os.Getenv("PATH")

func print_if_error_nil(output string, err error) {
	if err == nil {
		fmt.Print(string(output))
	} else {
		fmt.Fprintln(os.Stderr, "Error executing input:", err)
	}
}

func clean_string(str string) string {
	return str[:len(str)-1]
}

func path_is_valid(path string) bool {
	if _, search_err := os.Stat(path); search_err == nil {
		return true
	}
	return false
}

func process_command(command string) (string, string) {
	if strings.Count(command, " ") >= 1 {
		chunks := strings.SplitN(command, " ", 2)
		return chunks[0], strings.TrimSpace(chunks[1][:len(chunks[1])-1])
	}
	return strings.TrimSpace(command[:len(command)-1]), ""
}

func search_executable_path(exe_name string) string {
	dirs := strings.Split(PATH, string(filepath.ListSeparator))
	for i := 0; i < len(dirs); i++ {
		search_path := dirs[i] + string(os.PathSeparator) + exe_name
		if path_is_valid(search_path) {
			return search_path
		} else if matches, _ := filepath.Glob(search_path + ".*"); len(matches) > 0 {
			return matches[0]
		}
	}
	return ""
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
		case "pwd":
			directory, err := os.Getwd()
			print_if_error_nil(directory, err)
			fmt.Println()
		case "cd":
			if path_is_valid(arg) {
				os.Chdir(arg)
			} else {
				fmt.Println(command_keyword + ": " + arg + ": " + "No such file or directory")
			}
		case "type":
			if slices.Contains(builtin, arg) {
				fmt.Println(arg + " is a shell builtin")
			} else {
				search_result := search_executable_path(arg)
				if search_result == "" {
					fmt.Println(arg + ": not found")
				} else {
					fmt.Println(arg + " is " + search_result)
				}
			}
		default:
			search_result := search_executable_path(command_keyword)
			if search_result == "" {
				fmt.Println(clean_string(full_command) + ": command not found")
			} else {
				//ToDo: handle multiple argments
				command_result := exec.Command(command_keyword, arg)
				output, err := command_result.Output()
				print_if_error_nil(string(output), err)
				// if err == nil {
				// 	fmt.Print(string(output))
				// } else {
				// 	fmt.Fprintln(os.Stderr, "Error executing input:", err)
				// }
			}
		}

		if err != nil {
			//panic(err)
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

	}
}
