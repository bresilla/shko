package shko

import (
	"flag"
	"fmt"
	"os"

	"github.com/bresilla/dirk"
	term "github.com/bresilla/shko/term"
)

var (
	termWidth, termHeight = term.Size()
	wrap                  = true
	shortcut              = 17
	dirASwitch            = true
	dirBSwitch            = false
	barSpace              = 0
	sideSpace             = 0
	textLen               = 15
	startDir, _           = dirk.MakeFile(os.Getenv("PWD"))
	currentDir            = startDir
	childrens             = currentDir.ListDir()
	drawlist              = childrens
	changeDir             = true
	number                = 0
	scroll                = 0
	foreward              = false
	backward              = false
	appfolder             = os.Getenv("HOME") + "/.config/shko"
	setfolder             = appfolder + "/settings"
	tempfolder            = appfolder + "/templates"
	confFile              = appfolder + "/config"
	dirFile               = setfolder + "/chdir"
	tabFile               = setfolder + "/tabdirs"
	scriptsFile           = setfolder + "/scripts"
	freqFile              = setfolder + "/frecency"
	memFile               = setfolder + "/memory"
	bulkFile              = setfolder + "/rename"
	markFile              = setfolder + "/makrs"
	fileD, _              = os.Create(dirFile)
	memory, _             = loadFromFile(memFile)
	frecency, _           = loadFromFile(freqFile)
	swichero, _           = loadFromFile(tabFile)
	scripts               = map[string]string{}
	bookmark              = map[string]string{}
	copySlice             dirk.Files
	showIcons             = true
	showChildren          = false
	showSize              = false
	showDate              = false
	showMode              = false
	showMime              = false
	center                = false
	greater               = true
	sortSize              = false
	sortDate              = false
	sortName              = false
	sortType              = false
	sortMode              = false
	sortChildren          = false
	showBar               = false
	homeDir, _            = dirk.MakeFile(os.Getenv("HOME"))
	tempDir, _            = dirk.MakeFile(tempfolder)
	dirA                  = homeDir
	dirB                  = tabDir(tabFile)
	dir1                  dirk.File
	dir2                  dirk.File
	dir3                  dirk.File
	dir4                  dirk.File
	dir5                  dirk.File
	dir6                  dirk.File
	dir7                  dirk.File
	dir8                  dirk.File
	dir9                  dirk.File
	dir0                  dirk.File
)

func Flags() {
	flag.BoolVar(&dirk.DiskUse, "d", false, "")
	flag.BoolVar(&center, "c", true, "")
	flag.BoolVar(&showChildren, "n", false, "")
	flag.BoolVar(&showSize, "s", false, "")
	flag.BoolVar(&showIcons, "i", true, "")
	flag.BoolVar(&showMode, "m", false, "")
	flag.BoolVar(&showDate, "t", false, "")
	flag.BoolVar(&showBar, "b", false, "")
	flag.IntVar(&shortcut, "short", 17, "")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CreateDir(dirName string) bool {
	src, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			panic(err)
		}
		return true
	}
	if src.Mode().IsRegular() {
		return false
	}
	return false
}

func Main() {
	CreateDir(appfolder)
	createTemplates(tempfolder)
	initializeBookmarks()
	initializeScriptlist()

	fmt.Print("\033[?25l")
	Flags()
	Loop(childrens)
	term.ClearAll()
	fmt.Print("\033[?25h")

	manageTabDir(currentDir.Path)

	if changeDir {
		fileD.WriteString(currentDir.Path)
		addToFrecency(currentDir)
	} else {
		fileD.WriteString(startDir.Path)
	}

	saveToFile(memory, memFile)
	saveToFile(frecency, freqFile)
	saveToFile(swichero, tabFile)
}

func entryConditions() {

}

func prepList(childrens *dirk.Files) (drawlist dirk.Files) {
	barSpace = 0
	sideSpace = 0
	if center {
		showChildren = false
		showSize = false
		showDate = false
		showMode = false
		dirk.DiskUse = false
		showBar = false
		showMime = false
	}
	if len(*childrens) > 0 {
		textLen = (*childrens)[len(*childrens)-1].MaxPath()
	}
	if showBar {
		barSpace = 2
	}
	termHeight = termHeight - barSpace
	foreward = false
	backward = false
	if len(*childrens) != 0 {
		for i := range *childrens {
			(*childrens)[i].Active = false
		}
		if number >= len(*childrens) {
			number = len(*childrens) - 1
		} else if number < 0 {
			number = 0
		}
		(*childrens)[number].Active = true
		drawlist = *childrens
	}
	if len(*childrens) > termHeight-1 {
		if number > termHeight/2 {
			foreward = true
			backward = false
		} else if number < termHeight/2-2 {
			backward = true
			foreward = false
		}
		if len(*childrens) < termHeight {
			scroll = 0
		}
		if scroll <= 0 {
			scroll = 0
			backward = false
		} else if scroll >= len(*childrens)+1-termHeight {
			scroll = len(*childrens) + 1 - termHeight
			foreward = false
		}
		drawlist = (*childrens)[0+scroll : termHeight-1-barSpace+scroll]
	}

	if center && termHeight > len(drawlist) {
		barSpace += termHeight/2 - (len(drawlist) / 2)
		sideSpace = termWidth/2 - textLen/2 - 5
	}
	return
}

func Loop(childrens dirk.Files) {
	for {
		termWidth, termHeight = term.Size()
		drawlist := prepList(&childrens)
		SelectInList(&number, &scroll, &drawlist, &childrens, &currentDir)
		ascii, keycode, _ := term.GetChar()
		if ascii == 13 || ascii == shortcut { //----------------------------	enter, SHORTCUT (quit + chdir)
			break
		} else if ascii == 113 { //-----------------------------------------	q (quit)
			changeDir = false
			break
		} else if ascii == 3 { // ------------------------------------------	Ctrl+c (quit)
			changeDir = false
			break
		} else if keycode == 38 || ascii == 107 { // -----------------------	up, k (previous)
			shkoUp(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if keycode == 40 || ascii == 106 { // -----------------------	down, j (next)
			shkoDown(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if keycode == 37 || ascii == 104 { // -----------------------	left, h (back)
			shkoLeft(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if keycode == 39 || ascii == 108 { // -----------------------	right, l (enter + open)
			shkoRight(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 47 { // -----------------------------------------	/ (match)
			shkoMatch(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 37 { // -----------------------------------------	% (find)
			shkoFind(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 27 { //	----------------------------------------	ESC (refresh)
			childrens = currentDir.ListDir()
		} else if ascii == 32 { //	----------------------------------------	SPACE (select)
			shkoSelect(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 45 { //	----------------------------------------	- (recurr)
			shkoRecurr(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 44 { //	----------------------------------------	, (files)
			shkoFiles(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 46 { //	----------------------------------------	. (hidden)
			shkoHidden(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 9 { //	----------------------------------------	TAB (switch)
			shkoSwitch(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 117 { // ----------------------------------------	u (union)
			shkoUnion(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 105 { // ----------------------------------------	i (indent)
			shkoIndent(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 111 { // ----------------------------------------	o (open)
			shkoOpen(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 100 { // ----------------------------------------	d (delete)
			shkoDelete(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 97 { //	----------------------------------------	x (archive)
			shkoArchive(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 121 { //	----------------------------------------	y (yank)
			shkoYank(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 112 { //	----------------------------------------	p (paste)
			shkoPaste(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 114 { //-----------------------------------------	r (rename)
			shkoRename(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 110 { //	----------------------------------------	n (new)
			shkoNew(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 115 { // ----------------------------------------	s (script)
			shkoScript(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 103 { // ----------------------------------------	g (go-to)
			shkoGoTo(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 98 { // -----------------------------------------	b (bookmarks)
			shkoBookIt(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 126 { //	----------------------------------------	~ (home)
			shkoHome(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 119 { //	----------------------------------------	w (tabs)
			shkoTabs(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 120 { // ----------------------------------------	x (menu)
			shkoMenu(&currentDir, &childrens, &drawlist, &number, &scroll)
		} else if ascii == 122 { // ----------------------------------------	z (test)
		}
	}
}
