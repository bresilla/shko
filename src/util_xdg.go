package spejt

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/hashstructure"
)

func StringDirToFile(path string) (file File) {
	file = CurrentDir(path)
	return
}

func StatDir(dir string) {
	f, err := os.Stat(dir)
	ErrorCheck(err)
	fmt.Println(f)
}

func CurrentDir(dir string) (file File) {
	f, err := os.Stat(dir)
	if err != nil {
		return file
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

	file.Path = dir
	file.Name = name
	file.Parent = parent
	file.Size = f.Size()
	file.Mode = f.Mode()
	file.ModTime = f.ModTime()
	file.IsDir = f.IsDir()

	if string(name[0]) == "." {
		file.Hidden = true
	}
	if ComputeHashes {
		var h uint64
		h, err = hashstructure.Hash(file, nil)
		if err != nil {
			return file
		}
		file.Hash = h
	}
	return
}
