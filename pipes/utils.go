package pipes

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

// getRandomTempFiles returns temp files matching the patterns of ioutils.GetTempFile
// but allows you to get many with the same name type
func getRandomTempFiles(patterns []string) (files []string, err error) {
	var i int
	const attempts = 10

	tmpDir := tempDir()
	if _, err = os.Stat(tmpDir); err != nil {
		return
	}

	for i = 0; i < attempts; i++ {
		files = make([]string, len(patterns))
		randName := nextRandom()
		for num, pattern := range patterns {
			var prefix, suffix string
			if pos := strings.LastIndex(pattern, "*"); pos != -1 {
				prefix, suffix = pattern[:pos], pattern[pos+1:]
			} else {
				prefix = pattern
			}
			files[num] = filepath.Join(tmpDir, prefix+randName+suffix)
		}
		if filesAreValid(files) {
			return
		}
	}
	return nil, errors.New("Ran out of max attemtps to find random pipe names")
}

func filesAreValid(files []string) bool {
	for _, file := range files {
		if len(file) == 0 {
			return false
		}

		// Note !os.IsNotExist is not the same as os.IsExist cause of
		// unhandled error state. We explicitly want no file and no errors
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return false
		}
	}
	return true
}
