// Package copy implements copying of files.
//
// Overview
//
// These functions are somewhat similar to what command `cp -r` does in linux.
// Use File(src, dst) to copy one file/directory to another file/directory,
// it will know if it is the same file and will not copy then.
// Use Files(srcs, dst) to copy multiple files/directories into a directory.
//
// source: https://github.com/rvi64/copy
// this version has some modifications e.g. if file exists, it adds a number after the name

package dirk

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func IfExists(name string) string {
	if _, err := os.Stat(name); err == nil {
		i := 1
		for {
			if _, err := os.Stat(name + strconv.Itoa(i)); err == nil {
				i++
			} else {
				break
			}
		}
		return name + strconv.Itoa(i)
	}
	return name
}

func cpFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	var out *os.File
	defer in.Close()

	dst = IfExists(dst)

	out, err = os.Create(dst)
	if err != nil {
		return err
	}

	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, si.Mode())
}

func cpDir(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	dst = IfExists(dst)

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = cpDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = cpFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// File copies src to dst like unix command cp does
// if dst is a directory copies src into dst
// if dst doesnt exist dst will be a copy of src
// if src and dst are files dst will be replaced by src
// if src is a directories but dst is a file returns an error
func Copy(src, dst string) error {
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcinfo.IsDir() {
		dstinfo, err := os.Stat(dst)
		if err == nil {
			if os.SameFile(srcinfo, dstinfo) {
				return fmt.Errorf("directory is itself: %s", dst)
			}

			dst += "/" + filepath.Base(src)
			dst = IfExists(dst)

			return cpDir(src, dst)
		}

		return cpDir(src, dst)
	}

	dstinfo, err := os.Stat(dst)
	if err == nil {
		if dstinfo.IsDir() {
			return cpFile(src, dst+"/"+filepath.Base(src))
		}

		if os.SameFile(srcinfo, dstinfo) {
			return nil
		}

		return cpFile(src, dst)
	}

	return cpFile(src, dst)
}
