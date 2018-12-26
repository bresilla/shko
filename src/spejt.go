package spejt

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	term "github.com/tj/go/term"
)

var (
	incFolder     = true
	incFiles      = true
	incHidden     = false
	wrap          = true
	shortcut      = 113
	duMode        = false
	currentDir, _ = makeFile(os.Getenv("PWD"))
	startDir, _   = makeFile(os.Getenv("PWD"))
	changeDir     = true
	_, termHeight = term.Size()
	number        = 0
	scroll        = 0
	appfolder     = "/tmp/spejt"
	dirFile       = appfolder + "/chdir"
	memFile       = appfolder + "/memory"
	confFile      = appfolder + "/config"
	file, _       = os.Create(dirFile)
	showIcons     = true
	showChildren  = false
	showSize      = false
	showDate      = false
	showMode      = false
	cmd           *exec.Cmd
)

func Loop() {
	fmt.Print("\033[?25l")
	createDirectory(appfolder)
	term.ClearAll()
	children, parent := ListDirs(currentDir)
	for {
		_, termHeight = term.Size()
		termHeight = termHeight - 1
		var foreward = false
		var backward = false
		subdirs := children
		if len(subdirs) > termHeight-1 {
			if number > termHeight/2 {
				foreward = true
				backward = false
			} else if number < termHeight/2-2 {
				backward = true
				foreward = false
			}
			if scroll <= 0 {
				scroll = 0
				backward = false
			} else if scroll >= len(children)+1-termHeight {
				scroll = len(children) + 1 - termHeight
				foreward = false
			}
			subdirs = subdirs[0+scroll : termHeight-1+scroll]
		}
		SelectInList(number, subdirs, children)
		ascii, keycode, _ := GetChar()
		if ascii == 13 || ascii == shortcut || keycode == shortcut {
			break
		} else if ascii == 27 || ascii == 3 {
			changeDir = false
			break
		} else if keycode == 38 { //up
			if backward {
				scroll--
			} else {
				number--
			}
			if number < 0 {
				if wrap {
					number = len(subdirs) - 1
					scroll = len(children) - 1
				} else {
					number = 0
				}
			}
			continue
		} else if keycode == 40 { //down
			if foreward {
				scroll++
			} else {
				number++
			}
			if number > len(subdirs)-1 {
				if wrap {
					number = 0
					scroll = 0
				} else {
					number = len(subdirs) - 1
				}
			}
			continue
		} else if keycode == 37 { //left
			backward = false
			foreward = false
			oldDir := currentDir
			currentDir, _ = makeFile(parent.Path)
			children, parent = ListDirs(currentDir)
			number, scroll = find(children, oldDir)
			continue
		} else if keycode == 39 { //right
			if len(subdirs) == 0 {
				continue
			}
			if subdirs[number].IsDir {
				oldDir := currentDir
				currentDir, _ = makeFile(subdirs[number].Path)
				children, parent = ListDirs(currentDir)
				addToMemory(oldDir, currentDir)
				number, scroll = findInMemory(currentDir, children)
			} else {
				OpenFile(subdirs[number])
				fmt.Print("\033[?25l")
			}
			backward = false
			foreward = false
			continue
		} else {
			if ascii == 32 {
				children, parent = ListDirs(currentDir)
				term.MoveTo(0, termHeight)
				Print(HighLight, Black, White, "leader")
				ascii, _, _ := GetChar()
				if ascii == 110 {
					showChildren = !showChildren
				} else if ascii == 102 {
					showMode = !showMode
				} else if ascii == 109 {
					showDate = !showDate
				} else if ascii == 115 {
					showSize = !showSize
				} else if ascii == 19 {
					duMode = !duMode
				} else if ascii == 105 {
					showIcons = !showIcons
				} else if ascii == 71 {
					number = len(subdirs) - 1
					scroll = len(children) - 1
				} else if ascii == 103 {
					number = 0
					scroll = 0
				} else {
					fmt.Print(" ")
					toPrint := "ascii: " + strconv.Itoa(ascii)
					Print(HighLight, Black, White, toPrint)
					GetChar()
				}
				continue
			} else if ascii == 44 {
				incFiles = !incFiles
				children, parent = ListDirs(currentDir)
			} else if ascii == 46 {
				incHidden = !incHidden
				children, parent = ListDirs(currentDir)
			} else if ascii == 35 {
				wrap = !wrap
			} else {
				continue
			}
		}
	}
	fmt.Print("\033[?25h")
	fmt.Println()
	saveMemoryToFile(memory)
	if changeDir {
		file.WriteString(currentDir.Path)
	} else {
		file.WriteString(startDir.Path)
	}
}
