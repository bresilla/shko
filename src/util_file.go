package shko

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/karrick/godirwalk"
	"github.com/mitchellh/hashstructure"
)

func ByteCountSI(b int64) string {
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

func ByteCountIEC(b int64) string {
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
	Path      string
	Name      string
	Parent    string
	Children  int
	Mime      string
	Extension string
	IsDir     bool
	Hidden    bool
	Size      int64
	Mode      os.FileMode
	ModTime   time.Time
	Hash      uint64 `hash:"ignore"`
	Other     Other
}
type Other struct {
	Number     int
	Selected   bool
	Active     bool
	ParentPath string
	HumanSize  string
	Deep       int
	NameLength int
	Icon       string
}

func MakeFile(dir string) (file File, err error) {
	f, err := os.Stat(dir)
	if err != nil {
		return
	}
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
		ModTime: f.ModTime(),
		IsDir:   f.IsDir(),
	}

	if f.IsDir() {
		if duMode {
			file.Size = gothrough(dir)
			file.Other.HumanSize = ByteCountIEC(file.Size)
		} else {
			file.Other.HumanSize = "0 B"
		}
		file.Extension = ""
		file.Mime = "folder/folder"
		file.Other.Icon = categoryicons["folder/folder"]
		children, _ := godirwalk.ReadDirnames(dir, nil)
		file.Children = len(children)
	} else {
		extension := path.Ext(dir)
		mime, _, _ := mimetype.DetectFile(dir)
		file.Other.HumanSize = ByteCountIEC(f.Size())
		file.Extension = extension
		file.Mime = mime
		file.Other.Icon = fileicons[extension]
		if file.Other.Icon == "" {
			file.Other.Icon = categoryicons["file/default"]
		}
	}
	if string(name[0]) == "." {
		file.Hidden = true
	}
	hash, err := hashstructure.Hash(file, nil)
	if err != nil {
		return file, err
	}
	file.Hash = hash
	file.Other.NameLength = len(file.Name)
	file.Other.ParentPath = parentPath
	return
}

func fileList(recurrent bool, dir File) (paths []File, err error) {
	paths = []File{}
	if recurrent {
		err = godirwalk.Walk(dir.Path, &godirwalk.Options{
			Callback: func(osPathname string, de *godirwalk.Dirent) (err error) {
				file, _ := MakeFile(osPathname)
				paths = append(paths, file)
				return nil
			},
			Unsorted:      true,
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

func chooseFile(incFolder, incFiles, incHidden, recurrent bool, dir File) (list []File) {
	files := []File{}
	folder := []File{}
	hidden := []File{}
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
		list = hidden
	} else {
		for _, f := range hidden {
			if !f.Hidden {
				list = append(list, f)
			}
		}
	}
	for i, _ := range list {
		list[i].Other.Number = i
	}
	return
}

func ListFiles(dir File) (files []File, parent File) {
	list := chooseFile(incFolder, incFiles, incHidden, recurrent, dir)
	parent, _ = MakeFile(path.Dir(dir.Path))
	for _, d := range list {
		files = append(files, d)
	}
	return
}

func createDirectory(dirName string) bool {
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
