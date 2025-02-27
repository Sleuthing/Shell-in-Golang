package main

import (
	"bufio"
	"bytes"
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
var stdout_redir = regexp.MustCompile(" > | 1> ")
var stderr_redir = regexp.MustCompile(" 2> ")
var original_stdout = os.Stdout
var original_stderr = os.Stderr

// print functions

// print output if no error exists otherwise print error
func get_output_or_err_message(output string, err error) {
	if err == nil {
		// return output
		fmt.Fprint(os.Stdout, string(output))
	} else {
		//return fmt.Sprintf("Error executing input: %s", err)
		fmt.Fprintln(os.Stderr, "Error executing input:", err)
	}
}

func get_no_such_file_or_directory_message(cmd_keywrd string, dir string) {
	fmt.Fprintln(os.Stderr, cmd_keywrd+": "+dir+": No such file or directory")
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

// helper functions

func clean_command_clause(str string) string {
	str = strings.TrimSpace(strings.Trim(str, "\n"))
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
	if stdout_redir.MatchString(arg) {
		arg := strings.Replace(arg, " 1> ", " > ", -1)
		chunks := strings.Split(arg, " > ")
		return chunks[0], chunks[1]
	}
	return arg, ""
}

func check_for_stderr_redir(arg string) (string, string) {
	if stderr_redir.MatchString(arg) {
		chunks := strings.Split(arg, " 2> ")
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
		command_sentence_without_stderr_redirection, stderr_output_path := check_for_stderr_redir(clean_command)
		command_sentence_without_stdout_redirection, stdout_output_path := check_for_stdout_redir(command_sentence_without_stderr_redirection)
		command_keyword, arg_clause := process_command(command_sentence_without_stdout_redirection)
		var stdout_file = get_output_file(stdout_output_path)
		var stderr_file = get_output_file(stderr_output_path)
		if stdout_file != nil {
			os.Stdout = stdout_file
			defer stdout_file.Close()
		}
		if stderr_file != nil {
			os.Stderr = stderr_file
			defer stdout_file.Close()
		}
		switch command_keyword {
		case "exit":
			//fmt.Println("logout")
			os.Exit(0)
		case "echo":
			// output_string = strings.Trim(arg_clause, "\"'") + "\n"
			fmt.Println(strings.Trim(arg_clause, "\"'"))
		case "pwd":
			directory, err := os.Getwd()
			// output_string = get_output_or_err_message(directory, err) + "\n"
			get_output_or_err_message(directory, err)
			fmt.Println()
		case "cd":
			if arg_clause == "~" {
				os.Chdir(HOME)
			} else if path_is_valid(arg_clause) {
				os.Chdir(arg_clause)
			} else {
				// output_string = get_no_such_file_or_directory_message(command_keyword, arg_clause)
				get_no_such_file_or_directory_message(command_keyword, arg_clause)
			}
		case "type":
			if is_builtin(arg_clause) {
				fmt.Println(arg_clause + " is a shell builtin")
				// output_string = arg_clause + " is a shell builtin" + "\n"
			} else {
				search_result := search_executable_path(arg_clause)
				if search_result == "" {
					fmt.Println(arg_clause + ": not found")
					// output_string = arg_clause + ": not found" + "\n"
				} else {
					fmt.Println(arg_clause + " is " + search_result)
					// output_string = arg_clause + " is " + search_result + "\n"
				}
			}
		default:
			search_result := search_executable_path(command_keyword)
			if search_result == "" {
				// output_string = clean_command + ": command not found" + "\n"
				fmt.Println(clean_command + ": command not found")
			} else {
				args := strings.Split(arg_clause, " ")

				command_result := exec.Command(command_keyword, args...)

				var stdout, stderr bytes.Buffer
				command_result.Stdout = &stdout
				command_result.Stderr = &stderr

				err := command_result.Run()

				if _, error_assertion_success := err.(*exec.ExitError); error_assertion_success {
					fmt.Fprint(os.Stdout, stdout.String())
					fmt.Fprint(os.Stderr, stderr.String())
				} else if err != nil {
					fmt.Printf("Failed to run command: %v\n", err)
				} else {
					fmt.Fprint(os.Stdout, stdout.String())
				}

			}

		}

		if err != nil {
			//panic(err)
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		os.Stdout = original_stdout
		os.Stderr = original_stderr
	}
}
