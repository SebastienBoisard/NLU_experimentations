package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// Word is
type Word struct {
	Label string
	Links map[string]int
}

var wordMap map[string]Word

func readFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines, scanner.Err()
}

func parseSentence(s string) {

	s = strings.Replace(s, ".", "", -1)
	s = strings.Replace(s, ",", "", -1)
	s = strings.Replace(s, "\"", "", -1)
	s = strings.Replace(s, "’", "", -1)
	s = strings.Replace(s, "?", "", -1)
	s = strings.Replace(s, "”", "", -1)
	s = strings.Replace(s, "“", "", -1)
	s = strings.Replace(s, ":", "", -1)
	s = strings.Replace(s, "—", "", -1)
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)

	// Split on space.
	result := strings.Split(s, " ")

	for _, word1 := range result {
		for _, word2 := range result {
			if word1 == word2 {
				continue
			}

			w2, ok := wordMap[word2]
			if ok == false {
				w2 = Word{Label: word2, Links: map[string]int{"test": 1}}
				wordMap[word2] = w2
			}

			w1, ok := wordMap[word1]
			if ok == false {
				w1 = Word{Label: word1, Links: make(map[string]int)}
				wordMap[word1] = w1
			}

			w1.Links[word2]++
		}
	}
}

func main() {

	args := os.Args

	if len(args) != 2 {
		fmt.Println("Error with parameters")
		fmt.Println("Syntax:", args[0], "<corpus_filename>")
		return
	}

	filename := args[1]
	lines, err := readFile(filename)
	if err != nil {
		log.Fatalf("Error while reading file '%s': %s", filename, err)
	}

	wordMap = make(map[string]Word)

	for _, sentence := range lines {
		parseSentence(sentence)
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(wordMap)
	if err != nil {
		fmt.Println("Error while encoding wordMap, err:", err)
		return
	}

	t := time.Now()
	databaseName := t.Format("20060102150405")

	err = ioutil.WriteFile("database_"+databaseName+".db", buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error while writing file, err:", err)
		return
	}
}
