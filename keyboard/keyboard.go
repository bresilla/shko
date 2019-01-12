package shko

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	t "github.com/pkg/term"
)

// Returns either an ascii code, or (if input is an arrow) a Javascript key code.
func GetChar() (ascii int, keyCode int, err error) {
	term, _ := t.Open("/dev/tty")
	t.RawMode(term)
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = term.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// Up
			keyCode = 38
		} else if bytes[2] == 66 {
			// Down
			keyCode = 40
		} else if bytes[2] == 67 {
			// Right
			keyCode = 39
		} else if bytes[2] == 68 {
			// Left
			keyCode = 37
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
	} else {
		// Two characters read??
	}
	term.Restore()
	term.Close()
	return
}

func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}
