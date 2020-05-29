# GoExamples
Contains list of small go programs created to learn and demonstrate Golang capabilities

# goroutine_gennum
Accept command line argument "number of random number to generate", create four goroutines to generate random numbers, use channel to receive the numbers and print them out to STDOUT. Once the number of random number is reached, notify the goroutines to exit using done channel.

```bash
sfoo@SamLenoX270:~/go/src/github.com/samuel-foo/goexamples/goroutine_gennum$ go run main.go 12
starting goroutine: 1
starting goroutine: 2
starting goroutine: 3
starting goroutine: 4
81 887 847 59 81 318 425 540 456 300 694 511
Generated enough numbers, notify goroutines to close

goroutine finishing...

goroutine finishing...

goroutine finishing...

goroutine finishing...
```

# goroutine_searchwordinfiles
Accept word to search and a list of files, for each file kicks off a goroutine to search for occurence of the word, once the word is found, each goroutine pushes 1 to numChan and the function main will keep a running total. Once each goroutine finishes searching the file, it will push 1 to doneChan, the main function will keep track of number of times done is received, if the number reaches the number of goroutines/files, main function will exit polling, reports back the word count and returns. 

Note only whole word search is supported.

```bash
sfoo@SamLenoX270:~/go/src/github.com/samuel-foo/goexamples/goroutine_searchwordinfiles$ go run main.go -word "him" -files "../testdata_searchfiles/1.txt,../testdata_searchfiles/2.txt,../testdata_searchfiles/3.txt,../testdata_searchfiles/4.txt,../testdata_searchfiles/5.txt,../testdata_searchfiles/6.txt"
Searching files for word him
files=[../testdata_searchfiles/1.txt ../testdata_searchfiles/2.txt ../testdata_searchfiles/3.txt ../testdata_searchfiles/4.txt ../testdata_searchfiles/5.txt ../testdata_searchfiles/6.txt]
EOF exiting....goroutine ../testdata_searchfiles/3.txt exiting readFile()
EOF exiting....goroutine ../testdata_searchfiles/1.txt exiting readFile()
EOF exiting....goroutine ../testdata_searchfiles/2.txt exiting readFile()
EOF exiting....goroutine ../testdata_searchfiles/4.txt exiting readFile()
EOF exiting....goroutine ../testdata_searchfiles/6.txt exiting readFile()
EOF exiting....goroutine ../testdata_searchfiles/5.txt exiting readFile()
found 4 'him' word
```

# goroutine_wordcountbuffchan
Similar to wc in unix, count the number of words in the list of files supplied. For each file kicks off a goroutine, once a token is found by the goroutine it sends 1 to pulseChan. The main function polls the pulseChan and increment a running total once 1 is received. Once each goroutine is finished it sends 1 to doneChan and exit, the main function keeps track of number of times 1 is received from doneChan. Once it reaches the number of files/goroutines it will report back the word count and exit the program.

```bash
sfoo@SamLenoX270:~/go/src/github.com/samuel-foo/goexamples/goroutine_wordcountbuffchan$ go run main.go -w -f "../testdata_searchfiles/1.txt,../testdata_searchfiles/2.txt"
Counting word: true
File list: ../testdata_searchfiles/1.txt, ../testdata_searchfiles/2.txt
Number of tokens: 200

sfoo@SamLenoX270:~/go/src/github.com/samuel-foo/goexamples/goroutine_wordcountbuffchan$ wc -w ../testdata_searchfiles/1.txt ../testdata_searchfiles/2.txt
 131 ../testdata_searchfiles/1.txt
  69 ../testdata_searchfiles/2.txt
 200 total
```

# goroutine_wordcountsharemem
Same as goroutine_wordcountbuffchan above, except it is implemented using sync.Mutex to guard access to shared memory.
```json
sfoo@SamLenoX270:~/go/src/github.com/samuel-foo/goexamples/goroutine_wordcountsharemem$ go run main.go -w -f "../testdata_searchfiles/1.txt,../testdata_searchfiles/2.txt"
Counting word: true
File list: ../testdata_searchfiles/1.txt, ../testdata_searchfiles/2.txt
Number of tokens: 200
```