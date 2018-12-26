package spejt

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
	HumanSize  string
	Deep       int
	NameLength int
	Icon       string
}

func makeFile(dir string) (file File, err error) {
	f, err := os.Stat(dir)
	if err != nil {
		return
	}
	parent := "/"
	name := "/"
	if dir != "/" {
		dir = path.Clean(dir)
		parent, name = path.Split(dir)
		parent = strings.TrimRight(parent, "/")
		_, parent = path.Split(parent)
		if parent == "" {
			parent = "/"
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
	return
}

func listRecourFiles(dir string) (files []File, err error) {
	files = []File{}
	err = godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) (err error) {
			file, _ := makeFile(osPathname)
			files = append(files, file)
			return nil
		},
		Unsorted:      true,
		ScratchBuffer: make([]byte, 64*1024),
	})
	return
}

func listCurrentFiles(dir string) (files []File, err error) {
	files = []File{}
	children, err := godirwalk.ReadDirnames(dir, nil)
	if err != nil {
		return
	}
	sort.Strings(children)
	for _, child := range children {
		osPathname := path.Join(dir + "/" + child)
		file, _ := makeFile(osPathname)
		files = append(files, file)
	}
	return
}

func ListChooseCurrent(incFolder, incFiles, incHidden bool, dir string) (list []File) {
	files := []File{}
	paths := []File{}
	list, err := listCurrentFiles(dir)
	if err != nil {
		return
	}
	for _, f := range list {
		if f.IsDir {
			paths = append(paths, f)
		} else {
			files = append(files, f)
		}
	}
	list = []File{}
	if incFolder {
		for _, d := range paths {
			if incHidden {
				list = append(list, d)
			} else {
				if d.Hidden == false {
					list = append(list, d)
				}
			}
		}
	}
	if incFiles {
		for _, f := range files {
			if incHidden {
				list = append(list, f)
			} else {
				if f.Hidden == false {
					list = append(list, f)
				}
			}
		}
	}
	return
}

func ListDirs(dir File) (files []File, parent File) {
	list := ListChooseCurrent(incFolder, incFiles, incHidden, dir.Path)
	parent, _ = makeFile(path.Dir(dir.Path))
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
