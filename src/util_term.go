package spejt

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func Tty_Size() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	ErrorCheck(err)
	stringify_out := string(out)
	tty_size := strings.Split(stringify_out, " ")

	re := regexp.MustCompile("[0-9]+")
	width_re := re.FindAllString(tty_size[1], -1)
	heigh_re := re.FindAllString(tty_size[0], -1)

	width, err := strconv.Atoi(width_re[0])
	heigh, err := strconv.Atoi(heigh_re[0])
	ErrorCheck(err)

	return width, heigh
}

func DashBorder() string {
	width, _ := Tty_Size()
	var toPrint string
	for n := 1; n <= width; n++ {
		toPrint = toPrint + "-"
	}
	return toPrint
}
func ErrorCheck(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
