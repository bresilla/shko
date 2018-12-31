package shko

import (
	"log"
	"os"
)

var template = map[string][]byte{
	"go": []byte(`package main

import (
	"fmt"
)

func main() {
	fmt.Println("HELLO WORLD")
}`),
	"bash":    []byte(`#!/usr/bin/env bash`),
	"default": []byte(``),
}

func makeTemplate(name string, bytes []byte) {
	if _, err := os.Stat(tempfolder + "/" + name); err == nil {
		log.Print("File Exists")
	} else {
		newFileName := tempfolder + "/" + name
		newFile, _ := os.Create(newFileName)
		newFile.Write(bytes)
		newFile.Close()
	}
}

func createTemplates() {
	for key, value := range template {
		makeTemplate(key, value)
	}
}
