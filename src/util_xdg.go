package spejt

import (
	"os"
	"path"
	"strings"

	"github.com/mitchellh/hashstructure"
)

func EnvDir(env string) (file File) {
	path := os.Getenv(env)
	file = CurrentDir(path)
	return
}

func CurrentDir(dir string) (file File) {
	f, err := os.Stat(dir)
	if err != nil {
		return file
	}
	parentPath, name := path.Split(dir)
	parentPath = strings.TrimRight(parentPath, "/")
	_, parent := path.Split(parentPath)
	if parent == "" {
		parent = "root"
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
