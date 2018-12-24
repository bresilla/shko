package spejt

import (
	"io/ioutil"
	"strings"
)

var memory, _ = loadMemoryFromFile()

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

func saveMemoryToFile(array []string) {
	jointMem := strings.Join(array, ",")
	ioutil.WriteFile(memFile, []byte(jointMem), 0666)
}

func loadMemoryFromFile() (memory []string, err error) {
	jointMem, err := ioutil.ReadFile(memFile)
	memory = strings.Split(string(jointMem), ",")
	return
}
