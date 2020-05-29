package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// FilesToSearch struct
type FilesToSearch struct {
	files []string
}

func (s *FilesToSearch) String() string {
	return strings.Join(s.files, ", ")
}

// Set values
func (s *FilesToSearch) Set(value string) error {
	fs := strings.Split(value, ",")
	for i, f := range fs {
		fs[i] = strings.TrimSpace(f)
	}

	s.files = fs
	return nil
}

func main() {
	var searchWord string
	flag.StringVar(&searchWord, "word", "", "Search word")

	files := FilesToSearch{}
	flag.Var(&files, "files", "List of comma separated files to search surrounded by double quotes")

	flag.Parse()

	if flag.NFlag() != 2 {
		fmt.Println("Please provide a search string and list of comma separated files")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Searching files for word %s\n", searchWord)
	fmt.Printf("files=%+v\n", files.files)

	numChan := make(chan int, 5)
	done := make(chan int, 1)
	jobSize := len(files.files)

	for _, f := range files.files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			fmt.Printf("File %s does not exist, skipping...\n", f)
			continue
		}

		go func(file string) {
			readFile(file, searchWord, numChan, done)
		}(f)
	}

	sum := 0
	record := 0
L:
	for {
		select {
		case n := <-numChan:
			if n == 1 {
				sum++
			}
		case <-done:
			record++
		default:
			if record == jobSize {
				break L
			}
		}
	}

	close(numChan)
	close(done)

	fmt.Printf("found %d '%s' word\n", sum, searchWord)
}

// readFile func
func readFile(file, searchWord string, numChan chan<- int, done chan<- int) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", file, err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanWords)
	for {
		if !s.Scan() {
			if err := s.Err(); err != nil {
				fmt.Printf("Error scanning file: %s", err)
				os.Exit(1)
			}
			fmt.Printf("EOF exiting....")
			break
		}

		t := s.Text()
		//fmt.Printf("File %s, found %s\n", file, t)
		if strings.ToLower(t) == strings.ToLower(searchWord) {
			numChan <- 1
		}
	}

	done <- 1
	fmt.Printf("goroutine %s exiting readFile()\n", file)
}
