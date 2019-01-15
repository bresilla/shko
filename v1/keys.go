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
		dirk.DiskUse = false
		statBar = false
		topBar = false
		showMime = false
	}
	foreward = false
	backward = false
}

func prepList(childrens dirk.Files) (drawlist dirk.Files) {
	entryConditions()
	drawlist = childrens
	termHeight = termHeight - topSpace
	if len(drawlist) > termHeight-1 {
		if number > termHeight/2 {
			foreward = true
			backward = false
		} else if number < termHeight/2-2 {
			backward = true
			foreward = false
		}
		if len(childrens) < termHeight {
			scroll = 0
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
	if number > len(childrens) && len(childrens) != 0 {
		number = len(childrens)
	} else if number < 0 {
		number = 0
	}
	if len(childrens) != 0 {
		for i := range childrens {
			childrens[i].Active = false
		}
		childrens[number].Active = true
	}
	return
}

func fuzzyFind(childrens dirk.Files, currentDir dirk.File) (matched dirk.Files) {
	matched = childrens
	pattern := ""
	results := dirk.FindFrom(pattern, childrens)
	for {
		termWidth, termHeight = t.Size()
		drawlist = prepList(matched)
		SelectInList(number, scroll, drawlist, matched, currentDir)
		StatusWrite("Search for:")
		fmt.Print(pattern)
		ascii, keycode, _ := t.GetChar()
		runeString := string(rune(ascii))
		if ascii > 33 && ascii < 127 {
			pattern += runeString
		} else if ascii == 127 && len(pattern) > 0 {
			pattern = pattern[:len(pattern)-1]
		} else if ascii == 27 {
			matched = childrens
			break
		} else if ascii == 13 || keycode > 0 {
			break
		} else {
			continue
		}
		if pattern == "" {
			matched = childrens
		} else {
			matched = dirk.Files{}
			results = dirk.FindFrom(pattern, childrens)
			for _, r := range results {
				matched = append(matched, childrens[r.Index])
			}
		}
		number = 0
		scroll = 0
	}
	return
}

func Loop(childrens dirk.Files) {
	for {
		termWidth, termHeight = t.Size()
		drawlist := prepList(childrens)
		SelectInList(number, scroll, drawlist, childrens, currentDir)
		ascii, keycode, _ := t.GetChar()
		if ascii == 47 { // ------------------------------------------------	/
			childrens = fuzzyFind(childrens, currentDir)
		} else if keycode == 38 || ascii == 107 { // -----------------------	up
			if backward {
				scroll--
			} else {
				number--
			}
			if number < 0 {
				if wrap {
					number = len(drawlist) - 1
					scroll = len(childrens) - 1
				} else {
					number = 0
				}
			}
			continue
		} else if keycode == 40 || ascii == 106 { // -----------------------	down
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
		} else if keycode == 37 && !dirk.Recurrent || ascii == 104 { // ---------	left
			oldDir := currentDir
			currentDir, _ = dirk.MakeFile(currentDir.ParentPath)
			childrens = currentDir.ListDir()
			number, scroll = findFile(childrens, oldDir)
			backward = false
			foreward = false
			continue
		} else if keycode == 39 || ascii == 108 { // -----------------------	right
			if len(drawlist) == 0 {
				continue
			}
			if drawlist[number].IsDir {
				oldDir := currentDir
				currentDir, _ = dirk.MakeFile(drawlist[number].Path)
				childrens = currentDir.ListDir()
				addToMemory(oldDir, currentDir)
				number, scroll = findInMemory(currentDir, childrens)
			} else {
				currentDir.Select(childrens, number).Edit()
				fmt.Print("\033[?25l")
			}
			backward = false
			foreward = false
			continue
		} else if ascii == 13 || ascii == shortcut { //---------------------	enter + SHORTCUT
			t.ClearAll()
			break
		} else if ascii == 113 { //-----------------------------------------	q
			changeDir = false
			t.ClearAll()
			break
		} else if ascii == 3 { // ------------------------------------------	Ctrl+c
			changeDir = false
			t.ClearAll()
			break
		} else {
			if ascii == 32 { // --------------------------------------------	SPACE
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
				case 100: //	--------------------------------------------	z
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
				continue
			} else if ascii == 45 { //	-------------------------------------	- (recurr)
				dirk.Recurrent = !dirk.Recurrent
				dirk.IncFolder = !dirk.IncFolder
				dirk.IncHidden = false
				dirk.DiskUse = false
				childrens = currentDir.ListDir()
				if dirk.Recurrent {
					for i := range childrens {
						for j := range childrens[i].Ancestors {
							if len(childrens[i].Ancestors[j]) > termWidth/4 {
								childrens[i].Ancestors[j] = childrens[i].Ancestors[j][:termWidth/4] + "..."
							}
						}
						parent, _ := dirk.MakeFile(currentDir.ParentPath)
						childrens[i].Ancestors = childrens[i].Ancestors[parent.AncestorNr+1:]
						prepName := strings.Join(childrens[i].Ancestors, "/")
						childrens[i].Name = prepName
					}
				}
				number = 0
				scroll = 0
			} else if ascii == 44 { //	-------------------------------------	,
				dirk.IncFiles = !dirk.IncFiles
				childrens = currentDir.ListDir()
			} else if ascii == 46 { //	-------------------------------------	.
				dirk.IncHidden = !dirk.IncHidden
				childrens = currentDir.ListDir()
			} else if ascii == 35 { //	-------------------------------------	#
				wrap = !wrap
			} else if ascii == 9 { //	-------------------------------------	TAB
				if dirASwitch {
					if len(childrens) > 0 {
						dirA, _ = dirk.MakeFile(childrens[0].ParentPath)
					} else {
						dirA, _ = dirk.MakeFile(currentDir.ParentPath)
					}
					currentDir = dirB
					childrens = dirB.ListDir()
					number, scroll = findInMemory(currentDir, childrens)
					dirASwitch = false
					dirBSwitch = true
					showIcons = !showIcons
				} else {
					if len(childrens) > 0 {
						dirB, _ = dirk.MakeFile(childrens[0].ParentPath)
					} else {
						dirB, _ = dirk.MakeFile(currentDir.ParentPath)
					}
					currentDir = dirA
					childrens = dirA.ListDir()
					number, scroll = findInMemory(currentDir, childrens)
					dirBSwitch = false
					dirASwitch = true
					showIcons = !showIcons
				}
			} else if ascii == 122 { // ------------------------------------	z (test)
				currentDir.Select(childrens, number).Overite([]byte("TRIM" + "\n"))
				childrens = currentDir.ListDir()
			} else if ascii == 111 { // ------------------------------------	o (open)
				StatusWrite("Press \"o\" to OPEN or \"w\" to OPEN WITH...")
				ascii, _, _ = t.GetChar()
				switch ascii {
				case 111:
					currentDir.Select(childrens, number).Start()
				case 119:
					toOpenWith := StatusRead("Open with", "nvim")
					currentDir.Select(childrens, number).StartWith(toOpenWith)
				}
			} else if ascii == 100 && len(drawlist) > 0 { // ---------------	d (delete)
				StatusWrite("Press \"d\" to DELETE selected")
				ascii, _, _ = t.GetChar()
				if ascii == 100 {
					StatusWrite("Are you sure you want to delete? Y/N")
					ascii, _, _ = t.GetChar()
					if ascii == 121 || ascii == 89 {
						currentDir.Select(childrens, number).Delete()
					}
				}
				childrens = currentDir.ListDir()
				number--
			} else if ascii == 120 { //	------------------------------------	x (archive)
				StatusWrite("Press \"x\" to EXTRACT or \"a\" to ARCHIVE")
				ascii, _, _ = t.GetChar()
				if ascii == 120 {
					err := archiver.Unarchive(childrens[number].Path, childrens[number].Path+"_E")
					if err != nil {
						log.Fatal(err)
					}
				} else if ascii == 97 {
					archSlice := []string{}
					onList := false
					name := ""
					for i, file := range childrens {
						if childrens[i].Selected {
							archSlice = append(archSlice, file.Path)
							onList = true
							name = currentDir.Parent
						}
					}
					if !onList {
						archSlice = append(archSlice, childrens[number].Path)
						name = childrens[number].Name
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
				childrens = currentDir.ListDir()
			} else if ascii == 121 && len(drawlist) > 0 { //	------------	y (yank copy)
				StatusWrite("Press \"y\" to YANK selected")
				ascii, _, _ = t.GetChar()
				if ascii == 121 {
					copySlice = currentDir.Select(childrens, number)
					childrens = currentDir.ListDir()
				}
			} else if ascii == 112 { //	------------------------------------	p (paste copy)
				if len(copySlice) > 0 {
					StatusWrite("Press \"p\" to PASTE or \"m\" to MOVE")
					ascii, _, _ = t.GetChar()
					if ascii == 112 {
						copySlice.Paste(currentDir)
					} else if ascii == 109 {
						copySlice.Move(currentDir)
					}
					childrens = currentDir.ListDir()
				}
			} else if ascii == 114 && len(drawlist) > 0 { //----------------	r (rename)
				StatusWrite("Press \"r\" to RENAME selected")
				ascii, _, _ = t.GetChar()
				if ascii == 114 {
					selected := currentDir.Select(childrens, number)
					if len(selected) > 1 {
						selected.Rename()
					} else {
						newname := StatusRead("Rename "+childrens[number].Name+" to", childrens[number].Name)
						selected.Rename(newname)
					}
				}
				childrens = currentDir.ListDir()
			} else if ascii == 110 { //	-------------------------------------	n (new)
				StatusWrite("Press \"f\" to make new FILE, \"d\" to make new FOLDER or \"t\" to select from TEMPLATES")
				ascii, _, _ = t.GetChar()
				switch ascii {
				case 110, 102:
					name := StatusRead("Enter filename", "file")
					currentDir.Touch(name)
				case 100:
					name := StatusRead("Enter dirname", "dir")
					currentDir.Mkdir(name)
				case 116:
					number = 0
					scroll = 0
					childrens = tempDir.ListDir()
					for {
						termWidth, termHeight = t.Size()
						drawlist := prepList(childrens)
						SelectInList(number, scroll, drawlist, childrens, tempDir)
						ascii, keycode, _ := t.GetChar()
						if ascii == 13 { // ----	ENTER
							newFile, _ := dirk.MakeFiles([]string{drawlist[number].Path})
							newFile.Paste(currentDir)
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
									scroll = len(childrens) - 1
									if len(childrens) < termHeight {
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
				childrens = currentDir.ListDir()
			} else if ascii == 118 { //	-------------------------------------	v (select)
				drawlist[number].Selected = !drawlist[number].Selected
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
			} else if ascii == 115 { // ------------------------------------	s (script)
				StatusWrite("Press any key to launch script")
				ascii, _, _ = t.GetChar()
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
						scriptFiles, _ := dirk.MakeFiles([]string{scriptsFile})
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
						toRun := re.ReplaceAllString(script, childrens[number].Path)
						RunScript(toRun)
					}
				}
			} else if ascii == 103 { // ------------------------------------	g (go-to)
				name := StatusRead("Go-To", "")
				matched := matchFrecency(name)
				if _, err := os.Stat(matched); err == nil {
					currentDir, _ = dirk.MakeFile(matched)
					childrens = currentDir.ListDir()
				}
			} else if ascii == 98 { // -------------------------------------	b (bookmarks)
				StatusWrite("Press any key to go to the mark")
				ascii, _, _ = t.GetChar()
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
						markFiles, _ := dirk.MakeFiles([]string{markFile})
						markFiles.Edit()
						fmt.Print("\033[?25l")
					} else {
						addBookmark(ascii, currentDir.Path)
						saveBookmarks()
					}
				} else {
					bookdir, exists := readBookmarks(ascii)
					if exists {
						currentDir, _ = dirk.MakeFile(bookdir)
						childrens = currentDir.ListDir()
					}
				}
			} else if ascii == 126 { //	------------------------------------	~
				childrens = homeDir.ListDir()
			} else if ascii == 119 { //	------------------------------------	w (warps)
				StatusWrite("Pres SPACE then \"0\" to \"9\" keys to save warp")
				ascii, _, _ = t.GetChar()
				if ascii == 32 {
					ascii, _, _ = t.GetChar()
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
				}
			} else {
				switch ascii {
				case 48:
					if dir0.Path != "" {
						currentDir, _ = dirk.MakeFile(dir0.Path)
						childrens = currentDir.ListDir()
					} else {
						continue
					}
				case 49:
					if dir1.Path != "" {
						currentDir, _ = dirk.MakeFile(dir1.Path)
						childrens = currentDir.ListDir()
					} else {
						continue
					}
				case 50:
					if dir2.Path != "" {
						currentDir, _ = dirk.MakeFile(dir2.Path)
						childrens = currentDir.ListDir()
					} else {
						continue
					}
				case 51:
					if dir3.Path != "" {
						currentDir, _ = dirk.MakeFile(dir3.Path)
						childrens = currentDir.ListDir()
					} else {
						continue
					}
				case 52:
					if dir4.Path != "" {
						currentDir, _ = dirk.MakeFile(dir4.Path)
						childrens = currentDir.ListDir()
					} else {
						continue
					}
				case 53:
					if dir5.Path != "" {
						currentDir, _ = dirk.MakeFile(dir5.Path)
						childrens = currentDir.ListDir()
					} else {
						continue
					}
				case 54:
					if dir6.Path != "" {
						currentDir, _ = dirk.MakeFile(dir6.Path)
						childrens = currentDir.ListDir()
					} else {
						continue
					}
				case 55:
					if dir7.Path != "" {
						currentDir, _ = dirk.MakeFile(dir7.Path)
						childrens = currentDir.ListDir()
					} else {
						continue
					}
				case 56:
					if dir8.Path != "" {
						currentDir, _ = dirk.MakeFile(dir8.Path)
						childrens = currentDir.ListDir()
					} else {
						continue
					}
				case 57:
					if dir9.Path != "" {
						currentDir, _ = dirk.MakeFile(dir9.Path)
						childrens = currentDir.ListDir()
					} else {
						continue
					}
				}
				continue
			}
		}
	}
}