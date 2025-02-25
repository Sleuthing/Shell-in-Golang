package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var builtin = []string{"exit", "echo", "type", "pwd", "cd"}
var PATH = os.Getenv("PATH")
var HOME, _ = os.UserHomeDir()
var check_redir = regexp.MustCompile(" > | 1> ")
var original_stdout = os.Stdout

// print functions
func print_output_if_err_nil(output string, err error) {
	if err == nil {
		fmt.Print(string(output))
	} else {
		fmt.Fprintln(os.Stderr, "Error executing input:", err)
	}
}

func print_no_such_file_or_directory(cmd_keywrd string, dir string) {
	fmt.Println(cmd_keywrd + ": " + dir + ": " + "No such file or directory")
}

// boolean functions
func path_is_valid(path string) bool {
	if _, search_err := os.Stat(path); search_err == nil {
		return true
	}
	return false
}

func is_builtin(cmd_keywrd string) bool {
	return slices.Contains(builtin, cmd_keywrd)
}

func clean_command_clause(str string) string {
	fmt.Println(len(str))
	str = strings.TrimSpace(strings.Trim(str, "\n"))
	fmt.Println(len(str))
	return str
}

func process_command(command string) (string, string) {
	if strings.Count(command, " ") >= 1 {
		chunks := strings.SplitN(command, " ", 2)
		return chunks[0], strings.TrimSpace(chunks[1])
	}
	return strings.TrimSpace(command), ""
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

func check_for_stdout_redir(arg string) (string, string) {
	if check_redir.MatchString(arg) {
		arg := strings.Replace(arg, " 1> ", " > ", -1)
		chunks := strings.Split(arg, " > ")
		return chunks[0], chunks[1]
	}
	return arg, ""
}

func get_output_file(output_path string) *os.File {
	if output_path != "" {
		outfile, err := os.Create(output_path)

		if err != nil {
			panic(err)
		}

		return outfile
	}
	return nil
}

// all checks if all elements in the slice satisfy the condition defined by fn
// and returns problematic element if found
func all[T any](slice []T, fn func(T) bool) (bool, *T) {
	for _, item := range slice {
		if !fn(item) {
			return false, &item
		}
	}
	return true, nil
}

func main() {
	for i := 0; i < 100; i++ {
		fmt.Fprint(os.Stdout, "$ ")
		full_command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		clean_command := clean_command_clause(full_command)
		command_sentence, output_path := check_for_stdout_redir(clean_command)
		// fmt.Printf("Entire first part: %s \nOutput path: %s \n", command_sentence, output_path)
		command_keyword, arg_clause := process_command(command_sentence)
		// fmt.Printf("Command: %s \narg_clause: %s \n", command_keyword, arg_clause)
		var outfile = get_output_file(output_path)
		if outfile != nil {
			os.Stdout = outfile
			defer outfile.Close()
		}
		switch command_keyword {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Trim(arg_clause, "'"))
		case "pwd":
			directory, err := os.Getwd()
			print_output_if_err_nil(directory, err)
			fmt.Println()
		case "cd":
			if arg_clause == "~" {
				os.Chdir(HOME)
			} else if path_is_valid(arg_clause) {
				os.Chdir(arg_clause)
			} else {
				print_no_such_file_or_directory(command_keyword, arg_clause)
			}
		case "type":
			if is_builtin(arg_clause) {
				fmt.Println(arg_clause + " is a shell builtin")
			} else {
				search_result := search_executable_path(arg_clause)
				if search_result == "" {
					fmt.Println(arg_clause + ": not found")
				} else {
					fmt.Println(arg_clause + " is " + search_result)
				}
			}
		default:
			search_result := search_executable_path(command_keyword)
			if search_result == "" {
				fmt.Println(clean_command + ": command not found")
			} else {
				args := strings.Split(arg_clause, " ")
				safe_to_execute := true
				if command_keyword == "cat" {
					all_path_args_valid, invalid_path := all(args, path_is_valid)
					if !all_path_args_valid && invalid_path != nil {
						os.Stdout = original_stdout
						safe_to_execute = false
						print_no_such_file_or_directory(command_keyword, *invalid_path)
					}
				}
				if safe_to_execute {
					command_result := exec.Command(command_keyword, args...)
					output, err := command_result.Output()
					print_output_if_err_nil(string(output), err)
				}

			}

		}

		if err != nil {
			//panic(err)
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		os.Stdout = original_stdout
	}
}
