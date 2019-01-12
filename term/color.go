package shko

import (
	"fmt"
	"strconv"
)

type Color int
type Style int

const (
	None = Color(iota)
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Grey
)

const (
	Default   = Style(iota)
	HighLight = Style(1)
	Underline = Style(4)
	Flicker   = Style(5)
	AntiWhite = Style(7)
	Invisible = Style(8)
)

func ansiText2(style Style, fg Color, bg Color) string {
	if fg == None && bg == None {
		return ""
	}
	s := []byte("\x1b[")
	if style != Default {
		s = strconv.AppendUint(s, (uint64)(style-Default), 10)
	} else {
		s = strconv.AppendUint(s, (uint64)(Default), 10)
	}
	// fmt.Printf("s:%v\n", string(s))
	if bg != None {
		s = strconv.AppendUint(append(s, ";"...), 40+(uint64)(bg-Black), 10)
	}
	// fmt.Printf("s bg:%v\n", string(s))
	if fg != None {
		s = strconv.AppendUint(append(s, ";"...), 30+(uint64)(fg-Black), 10)
	}
	s = append(s, "m"...)
	// fmt.Printf("s fg:%v\n", string(s))
	return string(s)
}

func changeColorAndStyle(style Style, fg Color, bg Color) {
	if fg == None && bg == None {
		return
	}
	fmt.Print(ansiText2(style, fg, bg))
}

func SetStyle(style Style, fg Color, bg Color) {
	changeColorAndStyle(style, fg, bg)
}

func ResetStyle() {
	fmt.Print("\x1b[0m")
}

func Invert(active bool, style Style, color Color) {
	if active {
		SetStyle(style, Black, color)
	} else {
		SetStyle(style, color, None)
	}
}
