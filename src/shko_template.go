package shko

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var template = map[string][]byte{
	"go": []byte(`package main

import (
	"fmt"
)

func main() {
	fmt.Println("HELLO WORLD")
}`),
	"bash": []byte(`#!/bin/bash
STR="Hello World!"
echo $STR`),
	"assembly": []byte(`    global  _main
    extern  _printf

    section .text
_main:
    push    message
    call    _printf
    add     esp, 4
    ret
message:
    db  'Hello, World', 10, 0`),
	"c": []byte(`#include <stdio.h>

int main(void){
	printf("hello, world\n");
}`),
	"c++": []byte(`#include <iostream>

int main(){
	std::cout << "Hello, world!\n";
	return 0;
}`),
	"c#": []byte(`using System;

class Program{
	static void Main(string[] args){
		Console.WriteLine("Hello, world!");
	}
}`),
	"dart": []byte(`main() {
print('Hello World!');
}`),
	"delphi": []byte(`procedure TForm1.ShowAMessage;
begin
	ShowMessage('Hello World!');
end;`),
	"f#": []byte(`open System
Console.WriteLine("Hello World!")`),
	"haskell": []byte(`module Main where

main :: IO ()
main = putStrLn "Hello, World!"`),
	"java": []byte(`class HelloWorldApp {
	public static void main(String[] args) {
		System.out.println("Hello World!"); // Prints the string to the console.
	}
}`),
	"javascript": []byte(`console.log("Hello World!");`),
	"objective-c": []byte(`main(){
	puts("Hello World!");
	return 0;
}`),
	"rust": []byte(`use std::io;

fn main() {
	let mut line = String::new();
	io::stdin().read_line(&mut line).expect("reading stdin");
	
	let mut i: i64 = 0;
	for word in line.split_whitespace() {
		i += word.parse::<i64>().expect("trying to interpret your input as numbers");
	}
	println!("{}", i);
}`),
	"swift":    []byte(`println("Hello, world!")`),
	"default0": []byte(``),
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

func createTemplates(folder string) {
	createDirectory(folder)
	for key, value := range template {
		makeTemplate(key, value)
	}
}

var bookinit = map[string]string{
	"~": homeDir.Path,
	"d": homeDir.Path + "/Documents",
}

var bookmark = map[string]string{}

func ASCII(r []rune) int {
	switch {
	case 97 <= r[0] && r[0] <= 122:
		return int(r[0] - 32)
	case 65 <= r[0] && r[0] <= 90:
		return int(r[0] + 32)
	default:
		return int(r[0])
	}
}

func initializeBookmarks() {
	if _, err := os.Stat(markFile); err == nil {
		jointMem, err := ioutil.ReadFile(markFile)
		if err != nil {
			return
		}
		allBooks := strings.Split(string(jointMem), "\n")
		for _, el := range allBooks {
			arr := strings.Split(el, " > ")
			if arr[0] == "" || arr[1] == "" {
				continue
			}
			bookmark[arr[0]] = arr[1]
		}
	} else {
		newFile, _ := os.Create(markFile)
		for i, el := range bookinit {
			newFile.WriteString(i + " > " + el + "\n")
		}
		newFile.Close()
	}
}

func saveBookmarks() {
	newFile, _ := os.Create(markFile)
	for i, el := range bookmark {
		newFile.WriteString(i + " > " + el + "\n")
	}
	newFile.Close()
}

func addBookmark(ascii int, path string) {
	runeString := string(rune(ascii))
	bookmark[runeString] = path
}

func deleteBookmark(ascii int) {
	for i := range bookmark {
		runeInt := rune(i[0])
		if int(runeInt) == ascii {
			delete(bookmark, i)
		}
	}
	return
}

func readBookmarks(ascii int) (file string, exists bool) {
	for i, el := range bookmark {
		runeInt := rune(i[0])
		if int(runeInt) == ascii {
			if _, err := os.Stat(el); err == nil {
				file = el
				exists = true
			}
		}
	}
	return
}
