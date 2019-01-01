package shko

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

func OpenFile(file File) bool {
	var cmd *exec.Cmd
	editor := os.Getenv("EDITOR")

	if len(editor) > 0 {
		cmd = exec.Command(editor, file.Path)
	} else {
		if _, err := os.Stat("/usr/bin/sensible-editor"); err == nil {
			cmd = exec.Command("/usr/bin/sensible-editor", file.Path)
		} else {
			cmd = exec.Command("/usr/bin/env", "nvim", file.Path)
		}
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Println("Error:", err)
		return false
	}
	return true
}

func EditFile(name string) bool {
	var cmd *exec.Cmd
	editor := os.Getenv("EDITOR")

	if len(editor) > 0 {
		cmd = exec.Command(editor, name)
	} else {
		if _, err := os.Stat("/usr/bin/sensible-editor"); err == nil {
			cmd = exec.Command("/usr/bin/sensible-editor", name)
		} else {
			cmd = exec.Command("/usr/bin/env", "nvim", name)
		}
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Println("Error:", err)
		return false
	}
	return true
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
