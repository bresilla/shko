package spejt

import (
	"fmt"
	"os"
	"os/exec"

	term "github.com/buger/goterm"
)

var (
	incFolder     = true
	incFiles      = true
	incHidden     = false
	wrap          = true
	shortcut      = 113
	currentDir, _ = makeFile(os.Getenv("PWD"))
	startDir, _   = makeFile(os.Getenv("PWD"))
	changeDir     = true
	number        = 0
	scroll        = 0
	outDir        = "/tmp/spejt"
	file, _       = os.Create(outDir)
	showIcons     = true
	showChildren  = false
	cmd           *exec.Cmd
)

func Loop() {
	fmt.Print("\033[?25l")
	term.Flush()
	term.Clear()
	children, parent := ListDirs(currentDir)
	for {
		var foreward = false
		var backward = false
		subdirs := children
		if len(subdirs) > term.Height()-1 {
			if number > term.Height()/2 {
				foreward = true
				backward = false
			} else if number < term.Height()/2-2 {
				backward = true
				foreward = false
			}
			if scroll <= 0 {
				scroll = 0
				backward = false
			} else if scroll >= len(children)+1-term.Height() {
				scroll = len(children) + 1 - term.Height()
				foreward = false
			}
			subdirs = subdirs[0+scroll : term.Height()-1+scroll]
		}
		SelectInList(number, subdirs)
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
				cmd = exec.Command("xdg-open", subdirs[number].Path)
				fmt.Print("\033[?25l")
			}
			backward = false
			foreward = false
			continue
		} else {
			if ascii == 32 {
				children, parent = ListDirs(currentDir)
				ascii, keycode, _ := GetChar()
				if ascii == 110 {
					showChildren = !showChildren
				} else if ascii == 105 {
					showIcons = !showIcons
				} else if ascii == 71 {
					number = len(subdirs) - 1
					scroll = len(children) - 1
				} else if ascii == 103 {
					number = 0
					scroll = 0
				} else {
					println(ascii, "\t", keycode)
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
	if changeDir {
		file.WriteString(currentDir.Path)
	} else {
		file.WriteString(startDir.Path)
	}
}
