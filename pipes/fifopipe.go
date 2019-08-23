package pipes

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

type fifopipe struct {
	name   string
	flag   int
	file   *os.File
	opened bool
	closed bool
}

// NewFifoPipe returns an io.ReaadWriteCloser that maintains an in and an out
// pipe to write and read from for inter-process comunication
func NewFifoPipe(name string, flag int) (Pipe, error) {
	if err := unix.Mkfifo(name, 0600); err != nil {
		return nil, err
	}

	return &fifopipe{
		name:   name,
		flag:   flag,
		opened: false,
		closed: false,
	}, nil
}

// tryOpen is a private method to attempt an open on every read or write
func (pipe *fifopipe) tryOpen() (err error) {
	if !pipe.opened {
		err = pipe.Open()
	} else if pipe.closed {
		err = errors.New("Pipe has already been closed")
	}
	return
}

// Write will try and connect to a pipe if this is the first time or delegate to the file writer
func (pipe *fifopipe) Write(p []byte) (n int, err error) {
	if err = pipe.tryOpen(); err != nil {
		return
	}
	return pipe.file.Write(p)
}

// Read will try and connect to a pipe if this is the first time or delegate to the file reader
func (pipe *fifopipe) Read(p []byte) (n int, err error) {
	if err = pipe.tryOpen(); err != nil {
		return
	}
	return pipe.file.Read(p)
}

// Open will open the pipe if there is no file there already
func (pipe *fifopipe) Open() (err error) {
	if pipe.opened {
		return errors.New("Pipes can't be reopened")
	}
	pipe.file, err = os.OpenFile(pipe.name, pipe.flag, os.ModeNamedPipe)
	return
}

func (pipe *fifopipe) Name() StreamNames {
	return StreamNames{
		In:  pipe.name,
		Out: pipe.name,
	}
}

// Close will close all file connections and delete all the pipes
// an attempt is made for every opperation even if the last fails
func (pipe *fifopipe) Close() (err error) {
	err = pipe.file.Close()
	if e := os.Remove(pipe.name); e != nil {
		err = e
	}
	return
}
