package spejt

import (
	"log"
	"os"
	"os/exec"
)

func OpenFile(file File) bool {
	var cmd *exec.Cmd
	filepath := file.Path
	editor := os.Getenv("EDITOR")

	if len(editor) > 0 {
		cmd = exec.Command(editor, filepath)
	} else {
		if _, err := os.Stat("/usr/bin/sensible-editor"); err == nil {
			cmd = exec.Command("/usr/bin/sensible-editor", filepath)
		} else {
			cmd = exec.Command("/usr/bin/env", "nvim", filepath)
		}
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Println("Error:", err)
		log.Println("File not saved:", filepath)
		return false
	}
	return true
}
