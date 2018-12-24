package spejt

import (
	term "github.com/buger/goterm"
)

func find(list []File, actual File) (number, scroll int) {
	for i, el := range list {
		if el.Name == actual.Name {
			if i < term.Height()-1 {
				number = i
				scroll = 0
				break
			} else {
				number = 0
				scroll = i
			}
		} else {
			number = 0
			scroll = 0
		}
	}
	return
}
