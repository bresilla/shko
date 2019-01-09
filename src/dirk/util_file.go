package dirk

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/bresilla/godirwalk"
	"github.com/gabriel-vasile/mimetype"
)

var (
	IgnoreSlice = []string{}
	IgnoreRecur = []string{"node_modules", ".git"}
	IncFolder   = true
	IncFiles    = true
	IncHidden   = false
	Recurrent   = false
	DiskUse     = false
)

func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func byteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func gothrough(dir string) (size int64) {
	godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) (err error) {
			f, err := os.Stat(osPathname)
			if err != nil {
				return
			}
			size += f.Size()
			return nil
		},
		Unsorted:      true,
		ScratchBuffer: make([]byte, 64*1024),
	})
	return
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

var fileicons = map[string]string{
	".7z":       "",
	".ai":       "",
	".apk":      "",
	".avi":      "",
	".bat":      "",
	".bmp":      "",
	".bz2":      "",
	".c":        "",
	".c++":      "",
	".cab":      "",
	".cc":       "",
	".clj":      "",
	".cljc":     "",
	".cljs":     "",
	".coffee":   "",
	".conf":     "",
	".cp":       "",
	".cpio":     "",
	".cpp":      "",
	".css":      "",
	".cxx":      "",
	".d":        "",
	".dart":     "",
	".db":       "",
	".deb":      "",
	".diff":     "",
	".dump":     "",
	".edn":      "",
	".ejs":      "",
	".epub":     "",
	".erl":      "",
	".f#":       "",
	".fish":     "",
	".flac":     "",
	".flv":      "",
	".fs":       "",
	".fsi":      "",
	".fsscript": "",
	".fsx":      "",
	".gem":      "",
	".gif":      "",
	".go":       "",
	".gz":       "",
	".gzip":     "",
	".hbs":      "",
	".hrl":      "",
	".hs":       "",
	".htm":      "",
	".html":     "",
	".ico":      "",
	".ini":      "",
	".java":     "",
	".jl":       "",
	".jpeg":     "",
	".jpg":      "",
	".js":       "",
	".json":     "",
	".jsx":      "",
	".less":     "",
	".lha":      "",
	".lhs":      "",
	".log":      "",
	".lua":      "",
	".lzh":      "",
	".lzma":     "",
	".markdown": "",
	".md":       "",
	".mkv":      "",
	".ml":       "λ",
	".mli":      "λ",
	".mov":      "",
	".mp3":      "",
	".mp4":      "",
	".mpeg":     "",
	".mpg":      "",
	".mustache": "",
	".ogg":      "",
	".pdf":      "",
	".php":      "",
	".pl":       "",
	".pm":       "",
	".png":      "",
	".psb":      "",
	".psd":      "",
	".py":       "",
	".pyc":      "",
	".pyd":      "",
	".pyo":      "",
	".rar":      "",
	".rb":       "",
	".rc":       "",
	".rlib":     "",
	".rpm":      "",
	".rs":       "",
	".rss":      "",
	".scala":    "",
	".scss":     "",
	".sh":       "",
	".slim":     "",
	".sln":      "",
	".sql":      "",
	".styl":     "",
	".suo":      "",
	".t":        "",
	".tar":      "",
	".tgz":      "",
	".ts":       "",
	".twig":     "",
	".vim":      "",
	".vimrc":    "",
	".wav":      "",
	".xml":      "",
	".xul":      "",
	".xz":       "",
	".yml":      "",
	".zip":      "",
}

var categoryicons = map[string]string{
	"folder/folder": "",
	"file/default":  "",
}

type File struct {
	Number    int
	Path      string
	Name      string
	Parent    string
	Ancestors []string
	Childrens []string
	Mime      string
	Extension string
	IsDir     bool
	Hidden    bool
	Size      int64
	Mode      os.FileMode
	BrtTime   time.Time
	AccTime   time.Time
	ChgTime   time.Time
	Other     Other
}
type Other struct {
	Children   int
	Selected   bool
	Active     bool
	ParentPath string
	HumanSize  string
	Deep       int
	Ignore     bool
	Icon       string
}

type Files []File

func (e Files) String(i int) string    { return e[i].Name }
func (e Files) Len() int               { return len(e) }
func (e Files) Swap(i, j int)          { e[i], e[j] = e[j], e[i] }
func (e Files) Less(i, j int) bool     { return e[i].Name < e[j].Name }
func (e Files) SortSize(i, j int) bool { return e[i].Size < e[j].Size }
func (e Files) SortDate(i, j int) bool { return e[i].BrtTime.Before(e[j].BrtTime) }

func MakeFile(dir string) (file File, err error) {
	f, err := os.Stat(dir)
	if err != nil {
		return
	}
	osStat := f.Sys().(*syscall.Stat_t)
	parent := "/"
	parentPath := "/"
	name := "/"
	if dir != "/" {
		dir = path.Clean(dir)
		parentPath, name = path.Split(dir)
		parent = strings.TrimRight(parentPath, "/")
		_, parent = path.Split(parent)
		if parent == "" {
			parent = "/"
			parentPath = "/"
		}
	}
	file = File{
		Path:    dir,
		Name:    name,
		Parent:  parent,
		Size:    f.Size(),
		Mode:    f.Mode(),
		IsDir:   f.IsDir(),
		BrtTime: timespecToTime(osStat.Mtim),
		AccTime: timespecToTime(osStat.Atim),
		ChgTime: timespecToTime(osStat.Ctim),
	}

	if f.IsDir() {
		if DiskUse {
			file.Size = gothrough(dir)
			file.Other.HumanSize = byteCountIEC(file.Size)
		} else {
			file.Other.HumanSize = "0 B"
		}
		file.Extension = ""
		file.Mime = "folder/folder"
		file.Other.Icon = categoryicons["folder/folder"]
		file.Childrens, _ = godirwalk.ReadDirnames(dir, nil)
		file.Other.Children = len(file.Childrens)
	} else {
		extension := path.Ext(dir)
		mime, _, _ := mimetype.DetectFile(dir)
		file.Other.HumanSize = byteCountIEC(f.Size())
		file.Extension = extension
		file.Mime = mime
		file.Other.Icon = fileicons[extension]
		if file.Other.Icon == "" {
			file.Other.Icon = categoryicons["file/default"]
		}
	}
	file.Ancestors = strings.Split(dir, "/")
	file.Other.Deep = len(file.Ancestors)
	file.Other.ParentPath = parentPath
	if string(name[0]) == "." {
		file.Hidden = true
	}
	for _, s := range file.Ancestors {
		if s != "" && string(s[0]) == "." {
			file.Other.Ignore = true
			break
		}
	}
	return
}

func fileList(recurrent bool, dir File) (paths Files, err error) {
	paths = Files{}
	if recurrent {
		err = godirwalk.Walk(dir.Path, &godirwalk.Options{
			Callback: func(osPathname string, de *godirwalk.Dirent) (err error) {
				file, _ := MakeFile(osPathname)
				paths = append(paths, file)
				return nil
			},
			Unsorted:      true,
			NoHidden:      true,
			Ignore:        IgnoreRecur,
			ScratchBuffer: make([]byte, 64*1024),
		})
	} else {
		children, err := godirwalk.ReadDirnames(dir.Path, nil)
		if err != nil {
			return paths, err
		}
		sort.Strings(children)
		for _, child := range children {
			osPathname := path.Join(dir.Path + "/" + child)
			file, _ := MakeFile(osPathname)
			paths = append(paths, file)
		}
	}
	return
}

func chooseFile(incFolder, incFiles, incHidden, recurrent bool, dir File) (list Files) {
	files := Files{}
	folder := Files{}
	hidden := Files{}
	ignore := Files{}
	paths, _ := fileList(recurrent, dir)
	for _, f := range paths {
		if f.IsDir {
			folder = append(folder, f)
		} else {
			files = append(files, f)
		}
	}
	if incFolder {
		for _, d := range folder {
			hidden = append(hidden, d)
		}
	}
	if incFiles {
		for _, f := range files {
			hidden = append(hidden, f)
		}
	}
	if incHidden {
		ignore = hidden
	} else {
		for _, f := range hidden {
			if !f.Hidden {
				ignore = append(ignore, f)
			}
		}
	}
	if len(IgnoreSlice) > 0 {
		for _, f := range ignore {
			for _, s := range IgnoreSlice {
				if f.Name == s {
					break
				}
				list = append(list, f)
				break
			}
		}
	} else {
		list = ignore
	}
	for i, _ := range list {
		list[i].Number = i
	}
	return
}

func sortFiles(list Files, s bool) {
	if s {
		sort.Sort(Files(list))
	} else {
		sort.Sort(sort.Reverse(list))
	}
}

func ListFiles(dir File) (files Files, parent File) {
	list := chooseFile(IncFolder, IncFiles, IncHidden, Recurrent, dir)
	parent, _ = MakeFile(path.Dir(dir.Path))
	for _, d := range list {
		files = append(files, d)
	}
	return
}

func CreateDirectory(dirName string) bool {
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