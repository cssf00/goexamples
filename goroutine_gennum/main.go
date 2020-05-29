package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// ErrorLocation indicates which location the error occurs and what error is it
type ErrorLocation struct {
	location string
	err      error
}

const gorNum = 4

func main() {
	numCount, _ := strconv.Atoi(os.Args[1])

	done := make(chan bool, gorNum)
	numchan := make(chan int, gorNum)

	for i := 0; i < gorNum; i++ {
		fmt.Printf("starting goroutine: %d\n", i+1)
		go genRandNum(done, numchan)
	}

	actualCount := 0
L:
	for {
		select {
		case n := <-numchan:
			actualCount++
			if actualCount <= numCount {
				fmt.Printf("%d ", n)
			} else {
				fmt.Printf("\nGenerated enough numbers, notify goroutines to close\n")
				for i := 0; i < gorNum; i++ {
					done <- true
				}
				break L
			}
		}
	}

	close(numchan)
	close(done)
	time.Sleep(5 * time.Second)
}

func genRandNum(done <-chan bool, numchan chan<- int) {
	for {
		select {
		case <-done:
			fmt.Printf("\ngoroutine finishing...\n")
			return
		default:
			randNum := rand.Intn(1000)
			numchan <- randNum
			time.Sleep(500 * time.Millisecond)
		}
	}
}
