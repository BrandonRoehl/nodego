package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

const suffix = ".pipe"

// Random number state.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextRandom() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

func main() {
	var pipeFile string

	{
		var (
			err    error
			i      int
			tmpDir = os.TempDir()
		)
		if _, err = os.Stat(tmpDir); err != nil {
			panic(err)
		}
		for i = 0; i < 10; i++ {
			pipeFile = filepath.Join(tmpDir, nextRandom()+suffix)
			if _, err = os.Stat(pipeFile); os.IsNotExist(err) {
				break
			}
		}
		if i >= 10 {
			panic("ran out of attempts")
		}
	}

	fmt.Println(pipeFile)
	fmt.Println()

	if err := unix.Mkfifo(pipeFile, 0600); err != nil {
		panic(err)
	}

	// READ PIPE
	file, err := os.OpenFile(pipeFile, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}
	// Cleanup the pipe whenever it is closed
	defer func() {
		file.Close()
		if err := os.Remove(pipeFile); err != nil {
			panic(err)
		}
	}()

	// Read from the pipe until the pipe is closed
	reader := bufio.NewReader(file)
	dec := json.NewDecoder(reader)

	// Runtime loop to decode json objects into the interface
	for dec.More() {
		var i interface{}
		if err := dec.Decode(&i); err != nil {
			log.Fatal(err)
		}
		log.Println(&i, i)
	}

	// WRITE PIPE
	// f, err := os.OpenFile(pipeFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)

}
