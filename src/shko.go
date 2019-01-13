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
	statBar               = false
	topBar                = statBar
	topSpace              = 0
	sideSpace             = 0
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
	flag.BoolVar(&topBar, "b", false, "")
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
