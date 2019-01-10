package shko

import (
	"flag"
	"fmt"
	"os"

	. "./dirk"
	term "github.com/tj/go/term"
)

var (
	termWidth, termHeight = term.Size()
	colors                = true
	wrap                  = true
	shortcut              = 17
	dirASwitch            = true
	dirBSwitch            = false
	statBar               = false
	topBar                = statBar
	topSpace              = 0
	sideSpace             = 0
	startDir, _           = MakeFile(os.Getenv("PWD"))
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
	copySlice             Files
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
	homeDir, _            = MakeFile(os.Getenv("HOME"))
	tempDir, _            = MakeFile(tempfolder)
	dirA                  = homeDir
	dirB                  = tabDir(tabFile)
	dir1                  File
	dir2                  File
	dir3                  File
	dir4                  File
	dir5                  File
	dir6                  File
	dir7                  File
	dir8                  File
	dir9                  File
	dir0                  File
)

func Flags() {
	flag.BoolVar(&colors, "o", true, "disable colors - defualt enable")
	flag.BoolVar(&DiskUse, "d", false, "")
	flag.BoolVar(&center, "c", false, "")
	flag.BoolVar(&showChildren, "n", false, "")
	flag.BoolVar(&showSize, "s", false, "")
	flag.BoolVar(&showIcons, "i", true, "")
	flag.BoolVar(&showMode, "m", false, "")
	flag.BoolVar(&showDate, "t", false, "")
	flag.BoolVar(&topBar, "b", false, "")
	flag.IntVar(&shortcut, "short", 17, "")
	flag.Parse()

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
