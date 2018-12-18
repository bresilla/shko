package main

import (
	"fmt"
	"io/ioutil"

	"github.com/h2non/filetype"
)

func main() {
	buf, _ := ioutil.ReadFile("/home/bresilla/Sets/.wallpaper/angel_city_1080.jpg")
	match, _ := filetype.Match(buf)

	fmt.Println(match)

	if filetype.IsImage(buf) {
		fmt.Println("File is an image")
	} else {
		fmt.Println("Not an image")
	}
}
