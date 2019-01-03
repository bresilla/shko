package shko

import (
	"fmt"

	fuzzy "github.com/bresilla/fuzzy"
)

func fuzzyFind(childrens Files, currentDir File) (matched Files) {
	matched = childrens
	var pattern string
	for {
		drawlist = prepList(matched)
		SelectInList(number, scroll, drawlist, matched, currentDir)
		statusWrite("Search for:")
		fmt.Print(pattern)
		ascii, _, _ := GetChar()
		runeString := string(rune(ascii))
		if ascii > 33 && ascii < 127 {
			pattern += runeString
		} else if ascii == 127 && len(pattern) > 0 {
			pattern = pattern[:len(pattern)-1]
		} else if ascii == 27 {
			matched = childrens
			break
		} else {
			break
		}
		if pattern == "" {
			matched = childrens
		} else {
			matched = Files{}
			results := fuzzy.FindFrom(pattern, childrens)
			for _, r := range results {
				matched = append(matched, childrens[r.Index])
			}
		}
		number = 0
		scroll = 0
	}
	return
}
