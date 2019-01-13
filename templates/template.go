package shko

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
