package shko

import (
	"fmt"
	"os"
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

func prepList(childrens []File) (drawlist []File) {
	entryConditions()
	drawlist = childrens
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
		} else if scroll >= len(childrens)+1-termHeight {
			scroll = len(childrens) + 1 - termHeight
			foreward = false
		}
		drawlist = drawlist[0+scroll : termHeight-1-topSpace+scroll]
	}
	return
}

func Loop(childrens []File, parent File) {
	for {
		drawlist := prepList(childrens)
		SelectInList(number, scroll, drawlist, childrens, currentDir)
		ascii, keycode, _ := GetChar()
		if ascii == shortcut { //-------------------------------------------	SHORTCUT
			term.ClearAll()
			break
		} else if ascii == 3 || ascii == 113 { // --------------------------	q, Ctrl+c
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
					scroll = len(childrens) - 1
					if len(childrens) < termHeight {
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
			childrens, parent = ListFiles(currentDir)
			number, scroll = find(childrens, oldDir)
			continue
		} else if keycode == 39 { // ---------------------------------------	right
			if len(drawlist) == 0 {
				continue
			}
			if drawlist[number].IsDir {
				oldDir := currentDir
				currentDir, _ = MakeFile(drawlist[number].Path)
				childrens, parent = ListFiles(currentDir)
				addToMemory(oldDir, currentDir)
				number, scroll = findInMemory(currentDir, childrens)
			} else {
				OpenFile(drawlist[number])
				fmt.Print("\033[?25l")
			}
			backward = false
			foreward = false
			continue
		} else {
			if ascii == 32 { // --------------------------------------------	SPACE
				term.MoveTo(0, termHeight+1)
				Print(HighLight, Black, White, "leader")
				ascii, _, _ := GetChar()
				if ascii == 110 { //	------------------------------------	n
					showChildren = !showChildren
					center = false
				} else if ascii == 109 { //	--------------------------------	m
					showMode = !showMode
					center = false
				} else if ascii == 116 { //	--------------------------------	t
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
					scroll = len(childrens) - 1
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
				childrens, parent = ListFiles(currentDir)
				number = 0
				scroll = 0
			} else if ascii == 44 { //	-------------------------------------	,
				incFiles = !incFiles
				childrens, parent = ListFiles(currentDir)
			} else if ascii == 46 { //	-------------------------------------	.
				incHidden = !incHidden
				childrens, parent = ListFiles(currentDir)
			} else if ascii == 45 { //	-------------------------------------	-
				if dirASwitch {
					if len(childrens) > 0 {
						dirA, _ = MakeFile(childrens[0].Other.ParentPath)
					} else {
						dirA, _ = MakeFile(parent.Path)
					}
					currentDir = dirB
					childrens, parent = ListFiles(dirB)
					number, scroll = findInMemory(currentDir, childrens)
					dirASwitch = false
					dirBSwitch = true
				} else {
					if len(childrens) > 0 {
						dirB, _ = MakeFile(childrens[0].Other.ParentPath)
					} else {
						dirB, _ = MakeFile(parent.Path)
					}
					currentDir = dirA
					childrens, parent = ListFiles(dirA)
					number, scroll = findInMemory(currentDir, childrens)
					dirBSwitch = false
					dirASwitch = true
				}
			} else if ascii == 100 { //	-------------------------------------	d
				statusWrite("Press \"d\" to DELETE selected")
				ascii, _, _ = GetChar()
				if ascii == 100 {
					onList := false
					for i := range childrens {
						if childrens[i].Other.Selected {
							os.RemoveAll(childrens[i].Path)
							onList = true
						}
					}
					if !onList {
						os.RemoveAll(drawlist[number].Path)
					}
				}
				childrens, parent = ListFiles(currentDir)
			} else if ascii == 121 { //	-------------------------------------	y
				statusWrite("Press \"y\" to YANK selected")
				ascii, _, _ = GetChar()
				if ascii == 121 {
					copySlice = []File{}
					onList := false
					for i, file := range childrens {
						if childrens[i].Other.Selected {
							copySlice = append(copySlice, file)
							onList = true
						}
					}
					if !onList {
						copySlice = append(copySlice, drawlist[number])
					}
					childrens, parent = ListFiles(currentDir)
				}
			} else if ascii == 112 { //	-------------------------------------	p
				if len(copySlice) > 0 {
					statusWrite("Press \"p\" to PASTE here")
					ascii, _, _ = GetChar()
					if ascii == 112 {
						for _, el := range copySlice {
							Copy(el.Path, currentDir.Path)
						}
					}
					childrens, parent = ListFiles(currentDir)
				}
			} else if ascii == 109 { //	-------------------------------------	m
				if len(copySlice) > 0 {
					statusWrite("Press \"m\" to MOVE here")
					ascii, _, _ = GetChar()
					if ascii == 109 {
						for _, el := range copySlice {
							Copy(el.Path, currentDir.Path)
							os.RemoveAll(el.Path)
							copySlice = []File{}
						}
					}
					childrens, parent = ListFiles(currentDir)
				}
			} else if ascii == 114 { //	-------------------------------------	r
				statusWrite("Press \"r\" to RENAME selected")
				ascii, _, _ = GetChar()
				if ascii == 114 {
					onList := false
					for i, file := range childrens {
						if childrens[i].Other.Selected {
							onList = true
							print(file.Name)
						}
					}
					if !onList {
						newname := statusRead("Rename "+childrens[number].Name+" to: ", childrens[number].Name)
						os.Rename(childrens[number].Path, childrens[number].Other.ParentPath+"/"+newname)
					}
					childrens, parent = ListFiles(currentDir)
				}
			} else if ascii == 110 { //	-------------------------------------	n
				statusWrite("Press \"n\" to make new FILE or \"f\" to make new FOLDER")
				ascii, _, _ = GetChar()
				if ascii == 110 {
					newFile, _ := os.Create(currentDir.Path + "/" + "newFile.txt")
					newFile.Close()
				} else if ascii == 102 {
					os.MkdirAll(currentDir.Path+"/"+"newFolder", 0777)
				}
				childrens, parent = ListFiles(currentDir)
			} else if ascii == 126 { //	-------------------------------------	~
				childrens, parent = ListFiles(homeDir)
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
