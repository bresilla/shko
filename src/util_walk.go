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
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
	Hash    uint64 `hash:"ignore"`
}

func ListFiles(dir string) (files []File, err error) {
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

func ListPathsNFiles(dir string) (paths []File, files []File) {
	files = []File{}
	paths = []File{}
	list, err := ListFiles(dir)
	ErrorCheck(err)
	for _, f := range list {
		if f.IsDir {
			paths = append(paths, f)
		} else {
			files = append(files, f)
		}
	}
	return
}
