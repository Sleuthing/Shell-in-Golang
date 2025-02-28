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
	return strings.TrimSpace(strings.Trim(str, "\n"))
}

func remove_surrounding_quotes(str string) string {
	return strings.Trim(str, "\"'")
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

func check_for_redir(arg string, split_str string) (string, string) {
	var regex_match = regexp.MustCompile(split_str)
	if regex_match.MatchString(arg) {
		chunks := strings.Split(arg, split_str)
		return chunks[0], chunks[1]
	}
	return arg, ""
}

func get_output_file(output_path string, append bool) *os.File {
	var outfile *os.File
	var err error
	if output_path != "" {
		output_path = remove_surrounding_quotes(output_path)
		if append {
			outfile, err = os.OpenFile(output_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		} else {
			outfile, err = os.Create(output_path)
		}

		if err != nil {
			panic(err)
		}
		return outfile
	}
	return nil
}

func main() {
	for i := 0; i < 100; i++ {
		fmt.Fprint(os.Stdout, "$ ")
		full_command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		clean_command := clean_command_clause(full_command)
		clean_command = strings.Replace(clean_command, " 1> ", " > ", -1)
		clean_command = strings.Replace(clean_command, " 1>> ", " >> ", -1)
		command_sentence_without_stderr_append, stderr_append_output_path := check_for_redir(clean_command, " 2>> ")
		command_sentence_without_stdout_append, stdout_append_output_path := check_for_redir(command_sentence_without_stderr_append, " >> ")
		command_sentence_without_stderr_redirection, stderr_output_path := check_for_redir(command_sentence_without_stdout_append, " 2> ")
		command_sentence_without_stdout_redirection, stdout_output_path := check_for_redir(command_sentence_without_stderr_redirection, " > ")
		command_keyword, arg_clause := process_command(command_sentence_without_stdout_redirection)
		var stdout_file = get_output_file(stdout_output_path, false)
		var stdout_append_file = get_output_file(stdout_append_output_path, true)
		var stderr_file = get_output_file(stderr_output_path, false)
		var stderr_append_file = get_output_file(stderr_append_output_path, true)
		if stdout_file != nil {
			os.Stdout = stdout_file
			defer stdout_file.Close()
		}
		if stdout_append_file != nil {
			os.Stdout = stdout_append_file
			defer stdout_append_file.Close()
		}
		if stderr_file != nil {
			os.Stderr = stderr_file
			defer stderr_file.Close()
		}
		if stderr_append_file != nil {
			os.Stderr = stderr_append_file
			defer stderr_append_file.Close()
		}
		switch command_keyword {
		case "exit":
			//fmt.Println("logout")
			os.Exit(0)
		case "echo":
			arg_clause, _ := remove_surrounding_quotes(arg_clause), len(arg_clause)
			fmt.Println(strings.Join(strings.Fields(arg_clause), " "))
			// if old_arg_clause_length > len(arg_clause) {
			// 	fmt.Println(arg_clause)
			// } else {
			// 	fmt.Println(strings.Join(strings.Fields(arg_clause), " "))

			// }
			// fmt.Println()
		case "pwd":
			directory, err := os.Getwd()
			get_output_or_err_message(directory, err)
			fmt.Println()
		case "cd":
			if arg_clause == "~" {
				os.Chdir(HOME)
			} else if path_is_valid(arg_clause) {
				os.Chdir(arg_clause)
			} else {
				get_no_such_file_or_directory_message(command_keyword, arg_clause)
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
