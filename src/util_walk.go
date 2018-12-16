package spejt

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/karrick/godirwalk"
	"github.com/mitchellh/hashstructure"
)

var ComputeHashes = true

type File struct {
	Path    string
	Name    string
	Parent  string
	IsDir   bool
	Hidden  bool
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	Hash    uint64 `hash:"ignore"`
}

func ListRecourFiles(dir string) (files []File, err error) {
	files = []File{}
	err = godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) (err error) {
			f, err := os.Stat(osPathname)
			if err != nil {
				return
			}
			pathremain, name := path.Split(osPathname)
			pathremain = strings.TrimRight(pathremain, "/")
			_, parent := path.Split(pathremain)
			file := File{
				Path:    osPathname,
				Name:    name,
				Parent:  parent,
				Size:    f.Size(),
				Mode:    f.Mode(),
				ModTime: f.ModTime(),
				IsDir:   f.IsDir(),
			}
			if string(name[0]) == "." {
				file.Hidden = true
			}
			if ComputeHashes {
				var h uint64
				h, err = hashstructure.Hash(file, nil)
				if err != nil {
					return
				}
				file.Hash = h
			}
			files = append(files, file)
			return nil
		},
		Unsorted:      true,
		ScratchBuffer: make([]byte, 64*1024),
	})
	return
}

func ListRecourPathsNFiles(dir string) (paths []File, files []File) {
	files = []File{}
	paths = []File{}
	list, err := ListRecourFiles(dir)
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
	return
}

func ListCurrentFiles(dir string) (files []File, err error) {
	files = []File{}
	children, err := godirwalk.ReadDirnames(dir, nil)
	if err != nil {
		return
	}
	for _, child := range children {
		thepath := path.Join(dir + "/" + child)
		f, err := os.Stat(thepath)
		if err != nil {
			return files, err
		}
		_, parent := path.Split(dir)
		if parent == "" {
			parent = "/"
		}
		file := File{
			Path:    thepath,
			Name:    child,
			Parent:  parent,
			Size:    f.Size(),
			Mode:    f.Mode(),
			ModTime: f.ModTime(),
			IsDir:   f.IsDir(),
		}
		if string(child[0]) == "." {
			file.Hidden = true
		}
		if ComputeHashes {
			var h uint64
			h, err = hashstructure.Hash(file, nil)
			if err != nil {
				return files, err
			}
			file.Hash = h
		}
		files = append(files, file)
	}
	return
}

func ListCurrentPathsNFiles(dir string) (paths []File, files []File) {
	files = []File{}
	paths = []File{}
	list, err := ListCurrentFiles(dir)
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
	return
}

func ListChooseCurrent(incFolder, incFiles, incHidden bool, dir string) (list []File) {
	list = []File{}
	dirs, files := ListCurrentPathsNFiles(dir)
	if incFolder {
		for _, d := range dirs {
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
	parent = StringDirToFile(path.Dir(dir.Path))
	for _, d := range list {
		files = append(files, d)
	}
	return
}
