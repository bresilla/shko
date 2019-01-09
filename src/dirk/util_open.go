/*
	Open a file, directory, or URI using the OS's default
	application for that object type.  Optionally, you can
	specify an application to use.
	This is a proxy for the following commands:
	        OSX: "open"
	    	Windows: "start"
			Linux/Other: "xdg-open"
	Source: https://github.com/skratchdot/open-golang with modifications
*/

package dirk

import (
	"log"
	"os"
	"os/exec"
)

var cmd *exec.Cmd

/*
	Open a file, directory, or URI using the OS's default
	application for that object type. Wait for the open
	command to complete.
*/
func Run(input string) error {
	return open(input).Run()
}

/*
	Open a file, directory, or URI using the OS's default
	application for that object type. Don't wait for the
	open command to complete.
*/
func Start(input string) error {
	return open(input).Start()
}

/*
	Open a file, directory, or URI using the specified application.
	Wait for the open command to complete.
*/
func RunWith(input string, appName string) error {
	return openWith(input, appName).Run()
}

/*
	Open a file, directory, or URI using the specified application.
	Don't wait for the open command to complete.
*/
func StartWith(input string, appName string) error {
	return openWith(input, appName).Start()
}

func open(input string) *exec.Cmd {
	return exec.Command("xdg-open", input)
}

func openWith(input string, appName string) *exec.Cmd {
	return exec.Command(appName, input)
}

func Edit(file string) bool {
	editor := os.Getenv("EDITOR")
	if len(editor) > 0 {
		cmd = exec.Command(editor, file)
	} else {
		cmd = exec.Command("/usr/bin/env", "nvim", file)
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
