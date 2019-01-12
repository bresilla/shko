package shko

import (
	"io/ioutil"
	"math"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bresilla/shko/dirk"
)

var (
	frequency = 1
)

func tabDir(tabfile string) (file dirk.File) {
	file = homeDir
	if len(swichero) < 2 || len(swichero) > 2 {
		saveToFile([]string{homeDir.Path, homeDir.Path}, tabfile)
	} else {
		if _, err := os.Stat(swichero[0]); err == nil {
			file, _ = dirk.MakeFile(swichero[0])
		}
	}
	return
}

func manageTabDir(toadd string) {
	if toadd != swichero[1] {
		swichero = swichero[1:]
		swichero = append(swichero, toadd)
	}
}

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

func addToMemory(parent, child dirk.File) {
	for i, el := range memory {
		arr := strings.Split(el, " > ")
		if arr[0] == parent.Path {
			memory = memory[:i+copy(memory[i:], memory[i+1:])]
			break
		}
	}
	memory = append(memory, parent.Path+" > "+child.Path)
}

func findInMemory(parent dirk.File, child dirk.Files) (number, scroll int) {
	for _, el := range memory {
		arr := strings.Split(el, " > ")
		if arr[0] == parent.Path {
			file, _ := dirk.MakeFile(arr[1])
			number, scroll = findFile(child, file)
			break
		} else {
			number = 0
			scroll = 0
		}
	}
	return
}

func findFile(list dirk.Files, actual dirk.File) (number, scroll int) {
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

func addToFrecency(parent dirk.File) {
	for i, el := range frecency {
		arr := strings.Split(el, " > ")
		if len(arr) == 4 {
			if arr[2] == parent.Path {
				frequency, _ = strconv.Atoi(arr[0])
				frecency = frecency[:i+copy(frecency[i:], frecency[i+1:])]
				frequency++
			}
		}
	}
	timecal := int(calcFrecency(frequency, time.Now()))
	frecency = append(frecency, strconv.Itoa(frequency)+" > "+time.Now().String()+" > "+parent.Path+" > "+strconv.Itoa(timecal))
}

func calcFrecency(hits int, attime time.Time) (frecency float64) {
	toTime := float64(time.Now().Sub(attime))
	frecency = float64(hits)/0.25 + 3*math.Pow(10, -6)*toTime
	return
}

func matchFrecency(toMatch string) (matchedFile string) {
	var matchedList []string
	re := regexp.MustCompile(`(?i)` + toMatch)
	for _, el := range frecency {
		arr := strings.Split(el, " > ")
		_, name := path.Split(arr[2])
		if re.Match([]byte(name)) {
			matchedList = append(matchedList, el)
		}
	}
	if len(matchedList) > 0 {
		var bestScore float64
		for _, el := range matchedList {
			arr := strings.Split(el, " > ")
			frequency, _ = strconv.Atoi(arr[0])
			timeconv, _ := time.Parse("2019-01-01 19:07:28.623195367 +0100 CET m=+9.257016799", arr[1])
			if calcFrecency(frequency, timeconv) > bestScore {
				bestScore = calcFrecency(frequency, timeconv)
				matchedFile = arr[2]
			}
		}
	}
	return
}
