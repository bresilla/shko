package spejt

import (
	"os"
	"strings"
)

var memory = []string{"/:/home"}
var memFile, _ = os.Create("/tmp/spejt/memfile")

func addToMemory(parent, child File) {
	for i, el := range memory {
		arr := strings.Split(el, ":")
		if arr[0] == parent.Path {
			memory = memory[:i+copy(memory[i:], memory[i+1:])]
			break
		}
	}
	memory = append(memory, parent.Path+":"+child.Path)
}

func findInMemory(parent File, child []File) (number, scroll int) {
	for _, el := range memory {
		arr := strings.Split(el, ":")
		if arr[0] == parent.Path {
			file, _ := makeFile(arr[1])
			number, scroll = find(child, file)
			break
		} else {
			number = 0
			scroll = 0
		}
	}
	return
}
