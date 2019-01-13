package shko

import (
	"fmt"
	"log"

	"github.com/gobuffalo/packr/v2"
)

var Template = map[string][]byte{
	"go": []byte(`package main

import (
	"fmt"
)

func main() {
	fmt.Println("HELLO WORLD")
}`),
	"bash": []byte(`#!/bin/bash
STR="Hello World!"
echo $STR`)}

func Templates() {
	box := packr.New("myBox", "/home/bresilla/.go/src/github.com/bresilla/shko/templates/template_files")

	s, err := box.FindString("c.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}
