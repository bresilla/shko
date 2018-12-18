package spejt

func SelectInList(selected int, file []File) {
	for num, el := range file {
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
		Print(HighLight, Red, None, "\t »  "+file.Name+" /")
	} else if file.IsDir && !active {
		Print(HighLight, White, None, "\t»  "+file.Name+"/")
	} else if !file.IsDir && active {
		Print(HighLight, Red, None, "\t ♦  "+file.Name)
	} else if !file.IsDir && !active {
		Print(Default, Grey, None, "\t♦  "+file.Name)
	}
}
