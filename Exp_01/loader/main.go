package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

// Word is
type Word struct {
	Label string
	Links map[string]int
}

func main() {

	args := os.Args

	if len(args) != 2 {
		fmt.Println("Error with parameters")
		fmt.Println("Syntax:", args[0], "<database_filename>")
		return
	}

	filename := args[1]

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error while reading database, err:", err)
		return
	}

	wordMap := make(map[string]Word)

	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&wordMap)
	if err != nil {
		fmt.Println("Error while decoding wordMap, err:", err)
		return
	}

	var text string
	for {
		fmt.Print("\nEnter word: ")
		fmt.Scanln(&text)

		if text == "exit" {
			fmt.Println("bye")
			break
		}

		word, ok := wordMap[text]
		if ok == false {
			fmt.Println("word not found")
			continue
		}
		for linkLabel, linkCount := range word.Links {
			fmt.Println("    ", linkCount, "-", linkLabel)
		}
	}
}
