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
	appfolder         = "/tmp/shko"
	tempfolder        = appfolder + "/templates"
	dirFile           = appfolder + "/chdir"
	memFile           = appfolder + "/memory"
	confFile          = appfolder + "/config"
	bulkFile          = appfolder + "/rename"
	markFile          = appfolder + "/makrs"
	fileD, _          = os.Create(dirFile)
	copySlice         []File
	ignoreSlice       = []string{}
	showIcons         = true
	showChildren      = false
	showSize          = false
	showDate          = false
	showMode          = false
	duMode            = false
	center            = false
)

var (
	homeDir, _ = MakeFile(os.Getenv("HOME"))
	tempDir, _ = MakeFile(tempfolder)
	dirA       = homeDir
	dirB       = homeDir
	dir1       File
	dir2       File
	dir3       File
	dir4       File
	dir5       File
	dir6       File
	dir7       File
	dir8       File
	dir9       File
	dir0       File
)

func Flags() {
	flag.BoolVar(&duMode, "d", false, "")
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
	createDirectory(tempfolder)

	createTemplates()

	fmt.Print("\033[?25l")
	Flags()
	Loop(childrens, parent)
	fmt.Print("\033[?25h")

	saveMemoryToFile(memory)
	if changeDir {
		fileD.WriteString(currentDir.Path)
	} else {
		fileD.WriteString(startDir.Path)
	}
}
