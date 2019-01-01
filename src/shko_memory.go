package shko

import (
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"

	term "github.com/tj/go/term"
)

var (
	frequency int
)

func saveToFile(array []string, file string) {
	jointMem := strings.Join(array, "\n")
	ioutil.WriteFile(file, []byte(jointMem), 0666)
}

func loadFromFile(file string) (array []string, err error) {
	jointMem, err := ioutil.ReadFile(file)
	if err != nil {
		ioutil.WriteFile(file, []byte(""), 0666)
		return
	}
	array = strings.Split(string(jointMem), "\n")
	return
}

func addToMemory(parent, child File) {
	for i, el := range memory {
		arr := strings.Split(el, " > ")
		if arr[0] == parent.Path {
			memory = memory[:i+copy(memory[i:], memory[i+1:])]
			break
		}
	}
	memory = append(memory, parent.Path+" > "+child.Path)
}

func findInMemory(parent File, child []File) (number, scroll int) {
	for _, el := range memory {
		arr := strings.Split(el, " > ")
		if arr[0] == parent.Path {
			file, _ := MakeFile(arr[1])
			number, scroll = findFile(child, file)
			break
		} else {
			number = 0
			scroll = 0
		}
	}
	return
}

func findFile(list []File, actual File) (number, scroll int) {
	_, termHeight = term.Size()
	for i, el := range list {
		if el.Name == actual.Name {
			if i < termHeight/2 {
				number = i
				scroll = 0
				break
			} else {
				number = termHeight / 2
				scroll = el.Number - number
				break
			}
		} else {
			number = 0
			scroll = 0
		}
	}
	return
}

func findList(list []string, actual string) (number, scroll int) {
	_, termHeight = term.Size()
	for i, el := range list {
		if el == actual {
			if i < termHeight/2 {
				number = i
				scroll = 0
				break
			} else {
				number = termHeight / 2
				scroll = i - number
				break
			}
		} else {
			number = 0
			scroll = 0
		}
	}
	return
}

func addToFrecency(parent File) {
	for i, el := range frecency {
		arr := strings.Split(el, " > ")
		if arr[0] == parent.Path {
			frequency, _ = strconv.Atoi(arr[1])
			frecency = frecency[:i+copy(frecency[i:], frecency[i+1:])]
			frequency++
		}
	}
	frecency = append(frecency, parent.Path+" > "+strconv.Itoa(frequency)+" > "+time.Now().String())
}

func calcFrecency(hits int, attime time.Time) (frecency float64) {
	toTime := float64(time.Now().Sub(attime))
	frecency = float64(hits)/0.25 + 3*math.Pow(10, -6)*toTime
	return
}
