package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
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

	ct := NewCountTracker()
	var wg sync.WaitGroup
	for _, file := range fileList {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Fatalf("file %s does not exists\n", file)
		}
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			if err := scanToken(f, ct); err != nil {
				log.Printf("goroutine failed to scan file %s", f)
			}
		}(file)
	}
	wg.Wait()
	fmt.Printf("Number of tokens: %d\n", ct.tokenCount)
}

// NewCountTracker returns a new CountTracker
func NewCountTracker() *CountTracker {
	return &CountTracker{
		m: &sync.Mutex{},
	}
}

// CountTracker tracks the number token count
type CountTracker struct {
	m          *sync.Mutex
	tokenCount int
}

// IncrementCount increments the token count
func (t *CountTracker) IncrementCount() {
	t.m.Lock()
	defer t.m.Unlock()
	t.tokenCount++
}

func scanToken(file string, ct *CountTracker) error {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("Error opening file: %s\n", err)
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		ct.IncrementCount()
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error scanning next token: %s\n", err)
		return err
	}

	return nil
}
