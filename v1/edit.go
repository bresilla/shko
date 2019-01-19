package shko

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

var cmd *exec.Cmd

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

func RunScript(name string) bool {
	cmd = exec.Command("/bin/bash", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Println("Error:", err)
		return false
	}
	return true
}
