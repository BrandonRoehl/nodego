package nodego

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sys/unix"
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

func fileIsNotExist(file string) bool {
	_, err := os.Stat(file)
	return os.IsNotExist(err)
}

// getRandomPipeNames returns an in and out pipe and and error if one occured
func getRandomPipeNames() (pipeIn, pipeOut string, err error) {
	const (
		inSuffix, outSuffix = ".in.pipe", ".out.pipe"
		attempts            = 10 // Allowed collisions before we fail out
	)
	var i int
	tmpDir := os.TempDir()
	if _, err = os.Stat(tmpDir); err != nil {
		return
	}
	for i = 0; i < attempts; i++ {
		randName := nextRandom()
		pipeIn = filepath.Join(tmpDir, randName+inSuffix)
		pipeOut = filepath.Join(tmpDir, randName+outSuffix)
		if fileIsNotExist(pipeIn) && fileIsNotExist(pipeOut) {
			break
		}
	}
	if i >= attempts {
		pipeIn, pipeOut, err = "", "", errors.New("Ran out of max attemtps to find random pipe names")
	}
	return
}

// NewDuplexPipe returns an io.ReaadWriteCloser that maintains an in and an out
// pipe to write and read from for inter-process comunication
func NewDuplexPipe() (DuplexPipe, error) {
	pipeIn, pipeOut, err := getRandomPipeNames()
	if err != nil {
		return nil, err
	}

	fmt.Println(pipeIn)
	fmt.Println(pipeOut)
	fmt.Println()

	// READ PIPE
	inPipe, err := makePipe(pipeIn, os.O_CREATE)
	if err != nil {
		return nil, err
	}

	// WRITE PIPE
	outPipe, err := makePipe(pipeOut, os.O_WRONLY|os.O_CREATE|os.O_APPEND)
	if err != nil {
		return nil, err
	}

	return &duplexPipe{
		inPipe:  *inPipe,
		outPipe: *outPipe,
	}, nil
}

func makePipe(name string, flag int) (*pipe, error) {
	if err := unix.Mkfifo(name, 0600); err != nil {
		return nil, err
	}

	return &pipe{
		name:   name,
		flag:   flag,
		opened: false,
		closed: false,
	}, nil
}

// DuplexPipe is the interface to interact with the pipes that are created
type DuplexPipe interface {
	io.ReadWriteCloser
}

type duplexPipe struct {
	inPipe, outPipe pipe
}

func (pipe *duplexPipe) Write(p []byte) (n int, err error) {
	return pipe.outPipe.Write(p)
}

func (pipe *duplexPipe) Read(p []byte) (n int, err error) {
	return pipe.inPipe.Read(p)
}

func (pipe *duplexPipe) Close() (err error) {
	err = pipe.outPipe.Close()
	if e := pipe.inPipe.Close(); e != nil {
		err = e
	}
	return
}

type pipe struct {
	name   string
	flag   int
	file   *os.File
	opened bool
	closed bool
}

func (pipe *pipe) Write(p []byte) (n int, err error) {
	if err = pipe.tryOpen(); err != nil {
		return
	}
	return pipe.file.Write(p)
}

func (pipe *pipe) Read(p []byte) (n int, err error) {
	if err = pipe.tryOpen(); err != nil {
		return
	}
	return pipe.file.Read(p)
}

func (pipe *pipe) tryOpen() (err error) {
	if !pipe.opened {
		err = pipe.Open()
	} else if pipe.closed {
		err = errors.New("Pipe has already been closed")
	}
	return
}

func (pipe *pipe) Open() (err error) {
	if pipe.opened {
		return errors.New("Pipes can't be reopened")
	}
	pipe.file, err = os.OpenFile(pipe.name, pipe.flag, os.ModeNamedPipe)
	return
}

// Close will close all file connections and delete all the pipes
// an attempt is made for every opperation even if the last fails
func (pipe *pipe) Close() (err error) {
	err = pipe.file.Close()
	if e := os.Remove(pipe.name); e != nil {
		err = e
	}
	return
}