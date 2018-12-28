package spejt

import (
	"fmt"
	"strconv"

	term "github.com/tj/go/term"
)

func entryConditions() {
	topSpace = 0
	if topBar {
		topSpace = 2
		statBar = true
	} else {
		statBar = false
	}
	if center {
		showChildren = false
		showSize = false
		showDate = false
		showMode = false
		duMode = false
		statBar = false
		topBar = false
	}
	foreward = false
	backward = false
}

func prepList(filelist []File) (drawlist []File) {
	entryConditions()
	drawlist = filelist
	_, termHeight = term.Size()
	termHeight = termHeight - topSpace
	if len(drawlist) > termHeight-1 {
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
		} else if scroll >= len(filelist)+1-termHeight {
			scroll = len(filelist) + 1 - termHeight
			foreward = false
		}
		drawlist = drawlist[0+scroll : termHeight-1-topSpace+scroll]
	}
	return
}

func Loop(filelist []File, parent File) {
	for {
		drawlist := prepList(filelist)
		SelectInList(number, scroll, drawlist, filelist, currentDir)
		ascii, keycode, _ := GetChar()
		if ascii == 13 || ascii == shortcut { //----------------------------	q, SHORTCUT
			term.ClearAll()
			break
		} else if ascii == 3 { // ------------------------------------------	Ctrl+C
			changeDir = false
			term.ClearAll()
			break
		} else if ascii == 27 { // -----------------------------------------	ESC
			continue
		} else if keycode == 38 { // ---------------------------------------	up
			if backward {
				scroll--
			} else {
				number--
			}
			if number < 0 {
				if wrap {
					number = len(drawlist) - 1
					scroll = len(filelist) - 1
					if len(filelist) < termHeight {
						scroll = 0
					}
				} else {
					number = 0
				}
			}
			continue
		} else if keycode == 40 { // ---------------------------------------	down
			if foreward {
				scroll++
			} else {
				number++
			}
			if number > len(drawlist)-1 {
				if wrap {
					number = 0
					scroll = 0
				} else {
					number = len(drawlist) - 1
				}
			}
			continue
		} else if keycode == 37 { // ---------------------------------------	left
			backward = false
			foreward = false
			oldDir := currentDir
			currentDir, _ = MakeFile(parent.Path)
			filelist, parent = ListFiles(currentDir)
			number, scroll = find(filelist, oldDir)
			continue
		} else if keycode == 39 { // ---------------------------------------	right
			if len(drawlist) == 0 {
				continue
			}
			if drawlist[number].IsDir {
				oldDir := currentDir
				currentDir, _ = MakeFile(drawlist[number].Path)
				filelist, parent = ListFiles(currentDir)
				addToMemory(oldDir, currentDir)
				number, scroll = findInMemory(currentDir, filelist)
			} else {
				OpenFile(drawlist[number])
				fmt.Print("\033[?25l")
			}
			backward = false
			foreward = false
			continue
		} else {
			if ascii == 32 {
				term.MoveTo(0, termHeight+1)
				Print(HighLight, Black, White, "leader")
				ascii, _, _ := GetChar()
				if ascii == 110 { //	------------------------------------	n
					showChildren = !showChildren
					center = false
				} else if ascii == 102 { //	--------------------------------	f
					showMode = !showMode
					center = false
				} else if ascii == 109 { //	--------------------------------	m
					showDate = !showDate
					center = false
				} else if ascii == 98 { //	--------------------------------	b
					topBar = !topBar
					statBar = !statBar
					center = false
				} else if ascii == 115 { //	--------------------------------	s
					showSize = !showSize
					center = false
				} else if ascii == 99 { //	--------------------------------	c
					center = !center
				} else if ascii == 100 { //	--------------------------------	d
					duMode = !duMode
					center = false
				} else if ascii == 105 { //	--------------------------------	i
					showIcons = !showIcons
				} else if ascii == 71 { // ---------------------------------	G
					number = len(drawlist) - 1
					scroll = len(filelist) - 1
				} else if ascii == 103 { // --------------------------------	g
					number = 0
					scroll = 0
				} else {
					term.MoveTo(8, termHeight+1)
					toPrint := "ascii: " + strconv.Itoa(ascii)
					Print(HighLight, Black, White, toPrint)
					GetChar()
				}
				continue
			} else if ascii == 9 { //	-------------------------------------	TAB
				recurrent = !recurrent
				incFolder = !incFolder
				incHidden = false
				duMode = false
				filelist, parent = ListFiles(currentDir)
			} else if ascii == 44 { //	-------------------------------------	,
				incFiles = !incFiles
				filelist, parent = ListFiles(currentDir)
			} else if ascii == 46 { //	-------------------------------------	.
				incHidden = !incHidden
				filelist, parent = ListFiles(currentDir)
			} else if ascii == 45 { //	-------------------------------------	-
				if dirASwitch {
					if len(filelist) > 0 {
						dirA, _ = MakeFile(filelist[0].Other.ParentPath)
					} else {
						dirA, _ = MakeFile(parent.Path)
					}
					currentDir = dirB
					filelist, parent = ListFiles(dirB)
					number, scroll = findInMemory(currentDir, filelist)
					dirASwitch = false
					dirBSwitch = true
				} else {
					if len(filelist) > 0 {
						dirB, _ = MakeFile(filelist[0].Other.ParentPath)
					} else {
						dirB, _ = MakeFile(parent.Path)
					}
					currentDir = dirA
					filelist, parent = ListFiles(dirA)
					number, scroll = findInMemory(currentDir, filelist)
					dirBSwitch = false
					dirASwitch = true
				}
			} else if ascii == 126 { //	-------------------------------------	~
				filelist, parent = ListFiles(homeDir)
			} else if ascii == 35 { //	-------------------------------------	#
				wrap = !wrap
			} else if ascii == 118 { //	-------------------------------------	v
				drawlist[number].Other.Selected = !drawlist[number].Other.Selected
				if foreward {
					scroll++
				} else {
					number++
				}
				if number > len(drawlist)-1 {
					if wrap {
						number = 0
						scroll = 0
					} else {
						number = len(drawlist) - 1
					}
				}
				continue
			} else {
				continue
			}
		}
	}
}
