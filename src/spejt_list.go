package spejt

import (
	"fmt"
)

func ListDirs(dir string) {
	list := ListChooseCurrent(true, true, false, dir)
	for _, d := range list {
		fmt.Println(d.Name, "\t")
	}
}
