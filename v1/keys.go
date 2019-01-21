package shko

import (
	"fmt"

	"github.com/bresilla/dirk"
	t "github.com/bresilla/shko/term"
)

func entryConditions() {
	topSpace = 0
	if topBar {
		topSpace = 2
		statBar = true
	} else {
		statBar = false
	}
	if center {
		showChildren = false
		showSize = false
		showDate = false
		showMode = false
		dirk.DiskUse = false
		statBar = false
		topBar = false
		showMime = false
	}
	foreward = false
	backward = false
}

func prepList(childrens dirk.Files) (drawlist dirk.Files) {
	entryConditions()
	drawlist = childrens
	termHeight = termHeight - topSpace
	if len(drawlist) > termHeight-1 {
		if number > termHeight/2 {
			foreward = true
			backward = false
		} else if number < termHeight/2-2 {
			backward = true
			foreward = false
		}
		if len(childrens) < termHeight {
			scroll = 0
		}
		if scroll <= 0 {
			scroll = 0
			backward = false
		} else if scroll >= len(childrens)+1-termHeight {
			scroll = len(childrens) + 1 - termHeight
			foreward = false
		}
		drawlist = drawlist[0+scroll : termHeight-1-topSpace+scroll]
	}
	if number > len(childrens) && len(childrens) != 0 {
		number = len(childrens)
	} else if number < 0 {
		number = 0
	}
	if len(childrens) != 0 {
		for i := range childrens {
			childrens[i].Active = false
		}
		childrens[number].Active = true
	}
	return
}

func Loop(childrens dirk.Files) {
	for {
		termWidth, termHeight = t.Size()
		drawlist := prepList(childrens)
		SelectInList(number, scroll, drawlist, childrens, currentDir)
		ascii, keycode, _ := t.GetChar()
		if ascii == 13 || ascii == shortcut { //----------------------------	enter + SHORTCUT
			t.ClearAll()
			break
		} else if ascii == 113 { //-----------------------------------------	q
			changeDir = false
			t.ClearAll()
			break
		} else if ascii == 3 { // ------------------------------------------	Ctrl+c
			changeDir = false
			t.ClearAll()
			break
		} else if keycode == 38 || ascii == 107 { // -----------------------	up
			if backward {
				scroll--
			} else {
				number--
			}
			if number < 0 {
				if wrap {
					number = len(drawlist) - 1
					scroll = len(childrens) - 1
				} else {
					number = 0
				}
			}
			continue
		} else if keycode == 40 || ascii == 106 { // -----------------------	down
			if foreward {
				scroll++
			} else {
				number++
			}
			if number > len(drawlist)-1 {
				if wrap {
					number = 0
					scroll = 0
				} else {
					number = len(drawlist) - 1
				}
			}
			continue
		} else if keycode == 37 && !dirk.Recurrent || ascii == 104 { // ----	left
			oldDir := currentDir
			currentDir, _ = dirk.MakeFile(currentDir.ParentPath)
			childrens = currentDir.ListDir()
			number, scroll = findFile(childrens, oldDir)
			backward = false
			foreward = false
			continue
		} else if keycode == 39 || ascii == 108 { // -----------------------	right
			if len(drawlist) == 0 {
				continue
			}
			if drawlist[number].IsDir {
				oldDir := currentDir
				currentDir, _ = dirk.MakeFile(drawlist[number].Path)
				childrens = currentDir.ListDir()
				addToMemory(oldDir, currentDir)
				number, scroll = findInMemory(currentDir, childrens)
			} else {
				currentDir.Select(childrens).Edit()
				fmt.Print("\033[?25l")
			}
			backward = false
			foreward = false
			continue
		} else if ascii == 47 { // -----------------------------------------	/ (match)
			shkoMatch(&currentDir, &childrens, &drawlist, &number, &scroll, 47)
		} else if ascii == 37 { // -----------------------------------------	% (find)
			shkoFind(&currentDir, &childrens, &drawlist, &number, &scroll, 37)
		} else if ascii == 27 { //	----------------------------------------	ESC (refresh)
			childrens = currentDir.ListDir()
		} else if ascii == 32 { //	----------------------------------------	SPACE (select)
			shkoSelect(&currentDir, &childrens, &drawlist, &number, &scroll, 32)
		} else if ascii == 45 { //	----------------------------------------	- (recurr)
			shkoRecurr(&currentDir, &childrens, &drawlist, &number, &scroll, 45)
		} else if ascii == 44 { //	----------------------------------------	,
			dirk.IncFiles = !dirk.IncFiles
			childrens = currentDir.ListDir()
		} else if ascii == 46 { //	----------------------------------------	.
			dirk.IncHidden = !dirk.IncHidden
			childrens = currentDir.ListDir()
		} else if ascii == 35 { //	----------------------------------------	#
			wrap = !wrap
		} else if ascii == 9 { //	----------------------------------------	TAB (switch)
			shkoSwitch(&currentDir, &childrens, &drawlist, &number, &scroll, 9)
		} else if ascii == 117 { // ----------------------------------------	u (union)
			shkoUnion(&currentDir, &childrens, &drawlist, &number, &scroll, 117)
		} else if ascii == 105 { // ----------------------------------------	i (indent)
			shkoIndent(&currentDir, &childrens, &drawlist, &number, &scroll, 105)
		} else if ascii == 111 { // ----------------------------------------	o (open)
			shkoOpen(&currentDir, &childrens, &drawlist, &number, &scroll, 111)
		} else if ascii == 100 && len(drawlist) > 0 { // -------------------	d (delete)
			shkoDelete(&currentDir, &childrens, &drawlist, &number, &scroll, 100)
		} else if ascii == 97 { //	----------------------------------------	x (archive)
			shkoArchive(&currentDir, &childrens, &drawlist, &number, &scroll, 97)
		} else if ascii == 121 && len(drawlist) > 0 { //	----------------	y (yank)
			shkoYank(&currentDir, &childrens, &drawlist, &number, &scroll, 121)
		} else if ascii == 112 { //	----------------------------------------	p (paste)
			shkoPaste(&currentDir, &childrens, &drawlist, &number, &scroll, 112)
		} else if ascii == 114 && len(drawlist) > 0 { //--------------------	r (rename)
			shkoRename(&currentDir, &childrens, &drawlist, &number, &scroll, 114)
		} else if ascii == 110 { //	----------------------------------------	n (new)
			shkoNew(&currentDir, &childrens, &drawlist, &number, &scroll, 110)
		} else if ascii == 115 { // ----------------------------------------	s (script)
			shkoScript(&currentDir, &childrens, &drawlist, &number, &scroll, 115)
		} else if ascii == 103 { // ----------------------------------------	g (go-to)
			shkoGoTo(&currentDir, &childrens, &drawlist, &number, &scroll, 103)
		} else if ascii == 98 { // -----------------------------------------	b (bookmarks)
			shkoBookIt(&currentDir, &childrens, &drawlist, &number, &scroll, 98)
		} else if ascii == 126 { //	----------------------------------------	~ (home)
			childrens = homeDir.ListDir()
		} else if ascii == 119 { //	----------------------------------------	w (warps)
			shkoWarp(&currentDir, &childrens, &drawlist, &number, &scroll, 119)
		} else if ascii == 120 { // ----------------------------------------	x (menu)
			shkoMenu(&currentDir, &childrens, &drawlist, &number, &scroll, 120)
		} else if ascii == 122 { // ----------------------------------------	z (test)
		}
	}
}
