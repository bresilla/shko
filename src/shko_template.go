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
