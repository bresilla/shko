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
		} else if keycode == 38 || ascii == 107 { // ------------------------	up
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
		} else if keycode == 40 || ascii == 106 { // ------------------------	down
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
		} else if keycode == 37 || ascii == 104 { // ------------------------	left
			backward = false
			foreward = false
			oldDir := currentDir
			currentDir, _ = MakeFile(parent.Path)
			childrens, parent = ListFiles(currentDir)
			number, scroll = find(childrens, oldDir)
			continue
		} else if keycode == 39 || ascii == 108 { // ------------------------	right
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
				switch ascii {
				case 110: //	------------------------------------	n
					showChildren = !showChildren
					center = false
				case 109: //	--------------------------------	m
					showMode = !showMode
					center = false
				case 116: //	--------------------------------	t
					showDate = !showDate
					center = false
				case 98: //	--------------------------------	b
					topBar = !topBar
					statBar = !statBar
					center = false
				case 115: //	--------------------------------	s
					showSize = !showSize
					center = false
				case 99: //	--------------------------------	c
					center = !center
				case 100: //	--------------------------------	d
					duMode = !duMode
					center = false
				case 105: //	--------------------------------	i
					showIcons = !showIcons
				default:
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
			} else if ascii == 35 { //	-------------------------------------	#
				wrap = !wrap
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
					name := statusRead("Enter filename: ", "file.txt")
					newFileName := currentDir.Path + "/" + name
					newFileName = IfExists(newFileName)
					newFile, _ := os.Create(newFileName)
					newFile.Close()
				} else if ascii == 102 {
					name := statusRead("Enter filename: ", "folder")
					newFolderName := currentDir.Path + "/" + name
					newFolderName = IfExists(newFolderName)
					os.MkdirAll(newFolderName, 0777)
				}
				childrens, parent = ListFiles(currentDir)
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
			} else if ascii == 71 { // -------------------------------------	G
				number = len(drawlist) - 1
				scroll = len(childrens) - 1
			} else if ascii == 103 { // ------------------------------------	g
				number = 0
				scroll = 0
			} else if ascii == 126 { //	-------------------------------------	~
				childrens, parent = ListFiles(homeDir)
			} else if ascii == 119 { //	-------------------------------------	b
				statusWrite("Pres one of \"0\" to \"9\" keys to save this as WARPMARK")
				ascii, _, _ = GetChar()
				switch ascii {
				case 48:
					dir0 = currentDir
				case 49:
					dir1 = currentDir
				case 50:
					dir2 = currentDir
				case 51:
					dir3 = currentDir
				case 52:
					dir4 = currentDir
				case 53:
					dir5 = currentDir
				case 54:
					dir6 = currentDir
				case 55:
					dir7 = currentDir
				case 56:
					dir8 = currentDir
				case 57:
					dir9 = currentDir
				}
			} else {
				switch ascii {
				case 48:
					if dir0.Path != "" {
						currentDir, _ = MakeFile(dir0.Path)
						childrens, parent = ListFiles(currentDir)
					} else {
						continue
					}
				case 49:
					if dir1.Path != "" {
						currentDir, _ = MakeFile(dir1.Path)
						childrens, parent = ListFiles(currentDir)
					} else {
						continue
					}
				case 50:
					if dir2.Path != "" {
						currentDir, _ = MakeFile(dir2.Path)
						childrens, parent = ListFiles(currentDir)
					} else {
						continue
					}
				case 51:
					if dir3.Path != "" {
						currentDir, _ = MakeFile(dir3.Path)
						childrens, parent = ListFiles(currentDir)
					} else {
						continue
					}
				case 52:
					if dir4.Path != "" {
						currentDir, _ = MakeFile(dir4.Path)
						childrens, parent = ListFiles(currentDir)
					} else {
						continue
					}
				case 53:
					if dir5.Path != "" {
						currentDir, _ = MakeFile(dir5.Path)
						childrens, parent = ListFiles(currentDir)
					} else {
						continue
					}
				case 54:
					if dir6.Path != "" {
						currentDir, _ = MakeFile(dir6.Path)
						childrens, parent = ListFiles(currentDir)
					} else {
						continue
					}
				case 55:
					if dir7.Path != "" {
						currentDir, _ = MakeFile(dir7.Path)
						childrens, parent = ListFiles(currentDir)
					} else {
						continue
					}
				case 56:
					if dir8.Path != "" {
						currentDir, _ = MakeFile(dir8.Path)
						childrens, parent = ListFiles(currentDir)
					} else {
						continue
					}
				case 57:
					if dir9.Path != "" {
						currentDir, _ = MakeFile(dir9.Path)
						childrens, parent = ListFiles(currentDir)
					} else {
						continue
					}
				}
				continue
			}
		}
	}
}
