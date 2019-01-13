package shko

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	isatty "github.com/mattn/go-isatty"
)

// The text is rendered to the current cursor position, and when
// cleared with an empty string retains this position as if no
// text has been rendered.
func Renderer() func(string) {
	var prev string

	return func(curr string) {
		// clear lines
		if prev != "" {
			for range lines(prev) {
				MoveUp(1)
				ClearLine()
			}
		}

		// print lines
		if curr != "" {
			for _, s := range lines(curr) {
				fmt.Printf("%s\n", s)
			}
		}

		prev = curr
	}
}

// lines returns the lines in the given string.
func lines(s string) []string {
	return strings.Split(s, "\n")
}

// strip regexp.
var strip = regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)

// Strip ansi escape sequences.
func Strip(s string) string {
	return strip.ReplaceAllString(s, "")
}

// Length of characters with ansi escape sequences stripped.
func Length(s string) (n int) {
	for range Strip(s) {
		n++
	}
	return
}

// CenterLine a line of text.
func CenterLine(s string) string {
	r := strings.Repeat
	w, h := Size()
	size := Length(s)
	xpad := int(math.Abs(float64((w - size) / 2)))
	ypad := h / 2
	return r("\n", ypad) + r(" ", xpad) + s + r("\n", ypad)
}

// Size returns the width and height.
func Height() int {
	_, height := Size()
	return height
}

func Width() int {
	width, _ := Size()
	return width
}

func Size() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
	stringify_out := string(out)
	tty_size := strings.Split(stringify_out, " ")
	re := regexp.MustCompile("[0-9]+")
	width_re := re.FindAllString(tty_size[1], -1)
	heigh_re := re.FindAllString(tty_size[0], -1)
	width, _ := strconv.Atoi(width_re[0])
	heigh, _ := strconv.Atoi(heigh_re[0])
	return width, heigh
}

// ClearAll clears the screen.
func ClearAll() {
	fmt.Printf("\033[2J")
	MoveTo(1, 1)
}

// ClearLine clears the entire line.
func ClearLine() {
	fmt.Printf("\033[2K")
}

// ClearLineEnd clears to the end of the line.
func ClearLineEnd() {
	fmt.Printf("\033[0K")
}

// ClearLineStart clears to the start of the line.
func ClearLineStart() {
	fmt.Printf("\033[1K")
}

// MoveTo moves the cursor to (x, y).
func MoveTo(x, y int) {
	fmt.Printf("\033[%d;%df", y, x)
}

// MoveDown moves the cursor to the beginning of n lines down.
func MoveDown(n int) {
	fmt.Printf("\033[%dE", n)
}

// MoveUp moves the cursor to the beginning of n lines up.
func MoveUp(n int) {
	fmt.Printf("\033[%dF", n)
}

// SaveCursorPosition saves the cursor position.
func SaveCursorPosition() {
	fmt.Printf("\033[s")
}

// RestoreCursorPosition saves the cursor position.
func RestoreCursorPosition() {
	fmt.Printf("\033[u")
}

// HideCursor hides the cursor.
func HideCursor() {
	fmt.Printf("\033[?25l")
}

// ShowCursor shows the cursor.
func ShowCursor() {
	fmt.Printf("\033[?25h")
}

// IsTerminal returns true if fd is a tty.
func IsTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd)
}
