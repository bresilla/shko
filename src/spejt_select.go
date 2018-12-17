package spejt

func SelectInList(selected int, file []File) {
	for num, el := range file {
		if num < selected {
			Print(Default, None, None, el.Name)
		} else if num == selected {
			Print(HighLight, Red, None, el.Name)
		} else if num > selected {
			Print(Default, None, None, el.Name)
		}
	}
}
