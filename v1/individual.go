package shko

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bresilla/dirk"
	t "github.com/bresilla/shko/term"
	"github.com/mholt/archiver"
)

func shkoMenu(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	t.MoveTo(0, termHeight+1)
	Print(t.HighLight, t.Black, t.White, "leader")
	ascii, _, _ := t.GetChar()
	switch ascii {
	case 110: //	--------------------------------------------	n
		showChildren = !showChildren
		center = false
	case 117: //	--------------------------------------------	u
		showMime = !showMime
		center = false
	case 109: //	--------------------------------------------	m
		showMode = !showMode
		center = false
	case 116: //	--------------------------------------------	t
		showDate = !showDate
		center = false
	case 98: //	------------------------------------------------	b
		topBar = !topBar
		statBar = !statBar
		center = false
	case 115: //	--------------------------------------------	s
		showSize = !showSize
		center = false
	case 99: //	------------------------------------------------	c
		center = !center
	case 100: //	--------------------------------------------	d
		dirk.DiskUse = true
		showSize = true
		topBar = true
		statBar = true
		showDate = true
		showMode = true
		showMime = true
		showChildren = true
		center = false
	case 122: //	--------------------------------------------	z
		center = false
		dirk.DiskUse = !dirk.DiskUse
	case 105: //	--------------------------------------------	i
		showIcons = !showIcons
	default:
		t.MoveTo(8, termHeight+1)
		toPrint := "ascii: " + strconv.Itoa(ascii)
		Print(t.HighLight, t.Black, t.White, toPrint)
		t.GetChar()
	}
}

func shkoRecurr(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	dirk.Recurrent = !dirk.Recurrent
	dirk.DiskUse = false
	*childrens = currentDir.ListDir()
	if dirk.Recurrent {
		for i := range *childrens {
			(*childrens)[i].Name = (*childrens)[i].Path
		}
	}
	*number = 0
	*scroll = 0
}
func shkoSelect(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	(*drawlist)[*number].Selected = !(*drawlist)[*number].Selected
	if foreward {
		*scroll++
	} else {
		*number++
	}
	if *number > len(*drawlist)-1 {
		if wrap {
			*number = 0
			*scroll = 0
		} else {
			*number = len(*drawlist) - 1
		}
	}
}

func shkoSwitch(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	if dirASwitch {
		if len(*childrens) > 0 {
			dirA, _ = dirk.MakeFile((*childrens)[0].ParentPath)
		} else {
			dirA, _ = dirk.MakeFile(currentDir.ParentPath)
		}
		*currentDir = dirB
		*childrens = dirB.ListDir()
		*number, *scroll = findInMemory(*currentDir, *childrens)
		dirASwitch = false
		dirBSwitch = true
		showIcons = !showIcons
	} else {
		if len(*childrens) > 0 {
			dirB, _ = dirk.MakeFile((*childrens)[0].ParentPath)
		} else {
			dirB, _ = dirk.MakeFile(currentDir.ParentPath)
		}
		*currentDir = dirA
		*childrens = dirA.ListDir()
		*number, *scroll = findInMemory(*currentDir, *childrens)
		dirBSwitch = false
		dirASwitch = true
		showIcons = !showIcons
	}
}
func shkoFind(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	text := StatusRead("Write string to search", "text")
	*childrens = childrens.Find(dirk.Finder{Text: text})
}
func shkoMatch(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	matched := *childrens
	pattern := ""
	results := dirk.FindFrom(pattern, childrens)
	for {
		termWidth, termHeight = t.Size()
		*drawlist = prepList(matched)
		SelectInList(*number, *scroll, *drawlist, matched, *currentDir)
		StatusWrite("Search for:")
		fmt.Print(pattern)
		ascii, keycode, _ := t.GetChar()
		runeString := string(rune(ascii))
		if ascii > 33 && ascii < 127 {
			pattern += runeString
		} else if ascii == 127 && len(pattern) > 0 {
			pattern = pattern[:len(pattern)-1]
		} else if ascii == 27 {
			matched = *childrens
			break
		} else if ascii == 13 || keycode > 0 {
			break
		} else {
			continue
		}
		if pattern == "" {
			matched = *childrens
		} else {
			matched = dirk.Files{}
			results = dirk.FindFrom(pattern, childrens)
			for _, r := range results {
				matched = append(matched, (*childrens)[r.Index])
			}
		}
		*number = 0
		*scroll = 0
	}
	*childrens = matched
}

func shkoNew(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	StatusWrite("Press \"f\" to make new FILE, \"d\" to make new FOLDER or \"t\" to select from TEMPLATES")
	ascii, _, _ := t.GetChar()
	switch ascii {
	case 78, 70:
		read := StatusRead("Enter filenames", "file1 file2")
		names := strings.Split(read, " ")
		for _, name := range names {
			currentDir.Touch(name)
		}
	case 110, 102:
		name := StatusRead("Enter filename", "file")
		currentDir.Touch(name)
	case 68:
		read := StatusRead("Enter dirnames", "dir1 dir2")
		names := strings.Split(read, " ")
		for _, name := range names {
			currentDir.Mkdir(name)
		}
	case 100:
		name := StatusRead("Enter dirname", "dir")
		currentDir.Mkdir(name)
	case 116:
		number := 0
		scroll := 0
		*childrens = tempDir.ListDir()
		for {
			termWidth, termHeight = t.Size()
			drawlist := prepList(*childrens)
			SelectInList(number, scroll, drawlist, *childrens, tempDir)
			ascii, keycode, _ := t.GetChar()
			if ascii == 13 { // ----	ENTER
				newFile, _ := dirk.MakeFiles(drawlist[number].Path)
				newFile.Paste(*currentDir)
				break
			} else if keycode == 38 || ascii == 107 { //up
				if backward {
					scroll--
				} else {
					number--
				}
				if number < 0 {
					if wrap {
						number = len(drawlist) - 1
						scroll = len(*childrens) - 1
						if len(*childrens) < termHeight {
							scroll = 0
						}
					} else {
						number = 0
					}
				}
				continue
			} else if keycode == 40 || ascii == 106 { //down
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
			} else {
				break
			}
		}
	}
	*childrens = currentDir.ListDir()
}

func shkoRename(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	StatusWrite("Press \"r\" to RENAME selected")
	ascii, _, _ := t.GetChar()
	if ascii == 114 {
		selected := currentDir.Select(*childrens)
		oldname := (*childrens)[*number].Name
		if len(selected) > 1 {
			selected.Rename()
		} else {
			newname := StatusRead("Rename file to", oldname)
			selected.Rename(newname)
		}
	}
	*childrens = currentDir.ListDir()
}

func shkoPaste(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	if len(copySlice) > 0 {
		StatusWrite("Press \"p\" to PASTE or \"m\" to MOVE")
		ascii, _, _ := t.GetChar()
		if ascii == 112 {
			copySlice.Paste(*currentDir)
		} else if ascii == 109 {
			copySlice.Move(*currentDir)
		}
		*childrens = currentDir.ListDir()
	}
}

func shkoYank(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	StatusWrite("Press \"y\" to YANK selected")
	ascii, _, _ := t.GetChar()
	if ascii == 121 {
		copySlice = currentDir.Select(*childrens)
		*childrens = currentDir.ListDir()
	}
}

func shkoArchive(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	StatusWrite("Press \"x\" to EXTRACT or \"a\" to ARCHIVE")
	ascii, _, _ := t.GetChar()
	if ascii == 120 {
		err := archiver.Unarchive((*childrens)[*number].Path, (*childrens)[*number].Path+"_E")
		if err != nil {
			log.Fatal(err)
		}
	} else if ascii == 97 {
		archSlice := []string{}
		onList := false
		name := ""
		for i, file := range *childrens {
			if (*childrens)[i].Selected {
				archSlice = append(archSlice, file.Path)
				onList = true
				name = file.ParentPath
			}
		}
		if !onList {
			archSlice = append(archSlice, (*childrens)[*number].Path)
			name = (*childrens)[*number].Name
		}
		StatusWrite("Press \"t\" to TAR, \"z\" to ZIP or \"g\" to TGZ")
		ascii, _, _ = t.GetChar()
		if ascii == 116 {
			err := archiver.Archive(archSlice, name+".tar")
			if err != nil {
				log.Fatal(err)
			}
		} else if ascii == 122 {
			err := archiver.Archive(archSlice, name+".zip")
			if err != nil {
				log.Fatal(err)
			}
		} else if ascii == 103 {
			err := archiver.Archive(archSlice, name+".tar.gz")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	*childrens = currentDir.ListDir()
}

func shkoDelete(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	StatusWrite("Press \"d\" to DELETE selected")
	ascii, _, _ := t.GetChar()
	if ascii == 100 {
		StatusWrite("Are you sure you want to delete? Y/N")
		ascii, _, _ = t.GetChar()
		if ascii == 121 || ascii == 89 {
			currentDir.Select(*childrens).Delete()
		}
	}
	*childrens = currentDir.ListDir()
}

func shkoOpen(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	StatusWrite("Press \"o\" to OPEN or \"w\" to OPEN WITH...")
	ascii, _, _ := t.GetChar()
	switch ascii {
	case 111:
		currentDir.Select(*childrens).Start()
	case 119:
		toOpenWith := StatusRead("Open with", "nvim")
		currentDir.Select(*childrens).StartWith(toOpenWith)
	}
}

func shkoIndent(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	text := StatusRead("Enter name for indent directory", "dir")
	currentDir.Select(*childrens).Indent(text)
	*childrens = currentDir.ListDir()
}

func shkoUnion(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	text := StatusRead("Enter name for union entry", "entry")
	currentDir.Select(*childrens).Union(text)
	*childrens = currentDir.ListDir()
}

func shkoScript(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	StatusWrite("Press any key to launch script")
	ascii, _, _ := t.GetChar()
	if ascii == 32 {
		StatusWrite("Press any key to assign new script")
		ascii, _, _ = t.GetChar()
		_, exists := readScripts(ascii)
		if exists && ascii != 32 {
			runeString := string(rune(ascii))
			StatusWrite("Script exists, press \"" + runeString + "\" again to owerwrite")
			ascii2, _, _ := t.GetChar()
			if ascii2 == ascii {
				script := StatusRead("Write script", "file @")
				deleteScript(ascii)
				addScript(ascii, script)
				saveScript()
			}
		} else if ascii == 32 {
			scriptFiles, _ := dirk.MakeFiles(scriptsFile)
			scriptFiles.Edit()
			fmt.Print("\033[?25l")
		} else {
			script := StatusRead("Write script", "file @")
			addScript(ascii, script)
			saveScript()
		}
	} else {
		script, exists := readScripts(ascii)
		if exists {
			re := regexp.MustCompile(`@`)
			toRun := re.ReplaceAllString(script, (*childrens)[*number].Path)
			RunScript(toRun)
		}
	}
}

func shkoGoTo(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	name := StatusRead("Go-To", "")
	matched := matchFrecency(name)
	if _, err := os.Stat(matched); err == nil {
		*currentDir, _ = dirk.MakeFile(matched)
		*childrens = currentDir.ListDir()
	}
}

func shkoBookIt(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	StatusWrite("Press any key to go to the mark")
	ascii, _, _ := t.GetChar()
	if ascii == 32 {
		StatusWrite("Press any key to mark this directory")
		ascii, _, _ = t.GetChar()
		_, exists := readBookmarks(ascii)
		if exists && ascii != 32 {
			runeString := string(rune(ascii))
			StatusWrite("Mark exists, press \"" + runeString + "\" again to owerwrite")
			ascii2, _, _ := t.GetChar()
			if ascii2 == ascii {
				deleteBookmark(ascii)
				addBookmark(ascii, currentDir.Path)
				saveBookmarks()
			}
		} else if ascii == 32 {
			markFiles, _ := dirk.MakeFiles(markFile)
			markFiles.Edit()
			fmt.Print("\033[?25l")
		} else {
			addBookmark(ascii, currentDir.Path)
			saveBookmarks()
		}
	} else {
		bookdir, exists := readBookmarks(ascii)
		if exists {
			*currentDir, _ = dirk.MakeFile(bookdir)
			*childrens = currentDir.ListDir()
		}
	}
}

func shkoHome(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	*currentDir = homeDir
	*childrens = homeDir.ListDir()
}

func shkoWarp(currentDir *dirk.File, childrens, drawlist *dirk.Files, number, scroll *int, key int) {
	StatusWrite("Pres any key \"0\" to \"9\" ot go to tab")
	ascii, _, _ := t.GetChar()
	if ascii == 32 {
		ascii, _, _ = t.GetChar()
		switch ascii {
		case 48:
			dir0 = *currentDir
		case 49:
			dir1 = *currentDir
		case 50:
			dir2 = *currentDir
		case 51:
			dir3 = *currentDir
		case 52:
			dir4 = *currentDir
		case 53:
			dir5 = *currentDir
		case 54:
			dir6 = *currentDir
		case 55:
			dir7 = *currentDir
		case 56:
			dir8 = *currentDir
		case 57:
			dir9 = *currentDir
		}
	} else {
		switch ascii {
		case 48:
			if dir0.Path != "" {
				*currentDir, _ = dirk.MakeFile(dir0.Path)
				*childrens = currentDir.ListDir()
			}
		case 49:
			if dir1.Path != "" {
				*currentDir, _ = dirk.MakeFile(dir1.Path)
				*childrens = currentDir.ListDir()
			}
		case 50:
			if dir2.Path != "" {
				*currentDir, _ = dirk.MakeFile(dir2.Path)
				*childrens = currentDir.ListDir()
			}
		case 51:
			if dir3.Path != "" {
				*currentDir, _ = dirk.MakeFile(dir3.Path)
				*childrens = currentDir.ListDir()
			}
		case 52:
			if dir4.Path != "" {
				*currentDir, _ = dirk.MakeFile(dir4.Path)
				*childrens = currentDir.ListDir()
			}
		case 53:
			if dir5.Path != "" {
				*currentDir, _ = dirk.MakeFile(dir5.Path)
				*childrens = currentDir.ListDir()
			}
		case 54:
			if dir6.Path != "" {
				*currentDir, _ = dirk.MakeFile(dir6.Path)
				*childrens = currentDir.ListDir()
			}
		case 55:
			if dir7.Path != "" {
				*currentDir, _ = dirk.MakeFile(dir7.Path)
				*childrens = currentDir.ListDir()
			}
		case 56:
			if dir8.Path != "" {
				*currentDir, _ = dirk.MakeFile(dir8.Path)
				*childrens = currentDir.ListDir()
			}
		case 57:
			if dir9.Path != "" {
				*currentDir, _ = dirk.MakeFile(dir9.Path)
				*childrens = currentDir.ListDir()
			}
		}
	}
}
