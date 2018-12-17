package spejt

var (
	active     string
	deactive   string
	ttyW, ttyH = Tty_Size()
)

func SelectInList(selected int, file []File) {
	for num, el := range file {
		if el.IsDir {
			active = "\t »   " + el.Name + "/"
			deactive = "\t»   " + el.Name + "/"
		} else {
			active = "\t ≡   " + el.Name
			deactive = "\t≡   " + el.Name
		}
		if num < selected {
			colorList(el, false)
		} else if num == selected {
			colorList(el, true)
		} else if num > selected {
			colorList(el, false)
		}
	}
}

func colorList(file File, active bool) {
	if file.IsDir && active {
		Print(HighLight, Red, None, "\t ≡   "+file.Name)
	} else if file.IsDir && !active {
		Print(Default, None, None, "\t≡   "+file.Name)
	} else if !file.IsDir && active {
		Print(HighLight, Red, None, "\t »   "+file.Name)
	} else if !file.IsDir && !active {
		Print(Default, None, None, "\t»   "+file.Name)
	}
}
