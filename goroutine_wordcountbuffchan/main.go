package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// FileList is a list of file
type FileList []string

// String returns content of fileList
func (l *FileList) String() string {
	return strings.Join(*l, ", ")
}

// Set function parse the value into list of file names
func (l *FileList) Set(a string) error {
	for _, f := range strings.Split(a, ",") {
		*l = append(*l, strings.TrimSpace(f))
	}

	return nil
}

// Get function returns the value
func (l *FileList) Get() interface{} {
	return l
}

func main() {
	var (
		countWord bool
		fileList  FileList
	)

	flag.BoolVar(&countWord, "w", true, "Count number of words")
	flag.Var(&fileList, "f", "List of files comma separated")
	flag.Parse()

	fmt.Printf("Counting word: %t\n", countWord)
	fmt.Printf("File list: %s\n", fileList.Get())

	numOfFile := len(fileList)
	pulseChan := make(chan int, numOfFile)
	defer close(pulseChan)
	doneChan := make(chan int)
	defer close(doneChan)

	tokenCount := 0
	doneCount := 0

	for _, file := range fileList {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Fatalf("file %s does not exists\n", file)
		}
		go func(f string) {
			if err := scanToken(f, pulseChan, doneChan); err != nil {
				log.Printf("One goroutine failed to scan file %s", f)
			}
		}(file)
	}

L:
	for {
		select {
		case _ = <-pulseChan:
			tokenCount++
		case _ = <-doneChan:
			doneCount++
		default:
			if doneCount == numOfFile {
				// all go routine accounted for, exiting
				break L
			}
		}
	}

	fmt.Printf("Number of tokens: %d\n", tokenCount)
}

func scanToken(file string, pulse chan<- int, done chan<- int) error {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("Error opening file: %s\n", err)
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		pulse <- 1
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error scanning next token: %s\n", err)
		return err
	}

	done <- 1 // notify done
	return nil
}
