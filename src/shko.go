package shko

import (
	"flag"
	"fmt"
	"os"

	term "github.com/tj/go/term"
)

var (
	incFolder         = true
	incFiles          = true
	incHidden         = false
	recurrent         = false
	wrap              = true
	shortcut          = 13
	dirASwitch        = true
	dirBSwitch        = false
	statBar           = false
	topBar            = statBar
	topSpace          = 0
	sideSpace         = 0
	startDir, _       = MakeFile(os.Getenv("PWD"))
	currentDir        = startDir
	childrens, parent = ListFiles(currentDir)
	drawlist          = childrens
	changeDir         = true
	_, termHeight     = term.Size()
	number            = 0
	scroll            = 0
	foreward          = false
	backward          = false
	appfolder         = os.Getenv("HOME") + "/.config/shko"
	setfolder         = appfolder + "/settings"
	tempfolder        = appfolder + "/templates"
	confFile          = appfolder + "/config"
	dirFile           = setfolder + "/chdir"
	tabFile           = setfolder + "/tabdirs"
	scriptsFile       = setfolder + "/scripts"
	freqFile          = setfolder + "/frecency"
	memFile           = setfolder + "/memory"
	bulkFile          = setfolder + "/rename"
	markFile          = setfolder + "/makrs"
	fileD, _          = os.Create(dirFile)
	memory, _         = loadFromFile(memFile)
	frecency, _       = loadFromFile(freqFile)
	swichero, _       = loadFromFile(tabFile)
	scripts           = map[string]string{}
	bookmark          = map[string]string{}
	copySlice         Files
	ignoreSlice       = []string{}
	ignoreRecur       = []string{"node_modules", ".git"}
	showIcons         = true
	showChildren      = false
	showSize          = false
	showDate          = false
	showMode          = false
	showDu            = false
	showMime          = false
	center            = false
	greater           = true
	sortSize          = false
	sortDate          = false
	sortName          = false
	sortType          = false
	sortMode          = false
	sortChildren      = false
	homeDir, _        = MakeFile(os.Getenv("HOME"))
	tempDir, _        = MakeFile(tempfolder)
	dirA              = homeDir
	dirB              = tabDir(tabFile)
	dir1              File
	dir2              File
	dir3              File
	dir4              File
	dir5              File
	dir6              File
	dir7              File
	dir8              File
	dir9              File
	dir0              File
)

func Flags() {
	flag.BoolVar(&showDu, "d", false, "")
	flag.BoolVar(&center, "c", false, "")
	flag.BoolVar(&showChildren, "n", false, "")
	flag.BoolVar(&showSize, "s", false, "")
	flag.BoolVar(&showIcons, "i", true, "")
	flag.BoolVar(&showMode, "m", false, "")
	flag.BoolVar(&showDate, "t", false, "")
	flag.BoolVar(&topBar, "b", false, "")
	flag.Parse()

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Run() {
	createDirectory(appfolder)
	createTemplates(tempfolder)
	initializeBookmarks()
	initializeScriptlist()

	fmt.Print("\033[?25l")
	Flags()
	Loop(childrens, parent)
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
