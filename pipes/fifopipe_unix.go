// +build linux darwin

package pipes

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

// fifoPipe is a standard unix fifo pipe file wrapped in
// a convience contructor and io.ReadWriteCloser
type fifoPipe struct {
	name   string
	flag   int
	file   *os.File
	opened bool
	closed bool
}

// TempFifoPipe returns a pipe pointing at a temp file with O_RDWR O_APPEND privlages
func TempFifoPipe() (Pipe, error) {
	pipeNames, err := getRandomTempFiles([]string{"*.pipe"})
	if err != nil {
		return nil, err
	}
	return NewFifoPipe(pipeNames[0], os.O_RDWR|os.O_APPEND)
}

// NewFifoPipe returns an io.ReaadWriteCloser that maintains an in and an out
// pipe to write and read from for inter-process comunication
func NewFifoPipe(name string, flag int) (Pipe, error) {
	if err := unix.Mkfifo(name, 0600); err != nil {
		return nil, err
	}

	return &fifoPipe{
		name:   name,
		flag:   flag,
		opened: false,
		closed: false,
	}, nil
}

// tryOpen is a private method to attempt an open on every read or write
func (pipe *fifoPipe) tryOpen() (err error) {
	if !pipe.opened {
		err = pipe.Open()
	} else if pipe.closed {
		err = errors.New("Pipe has already been closed")
	}
	return
}

// Write will try and connect to a pipe if this is the first time or delegate to the file writer
func (pipe *fifoPipe) Write(p []byte) (n int, err error) {
	if err = pipe.tryOpen(); err != nil {
		return
	}
	return pipe.file.Write(p)
}

// Read will try and connect to a pipe if this is the first time or delegate to the file reader
func (pipe *fifoPipe) Read(p []byte) (n int, err error) {
	if err = pipe.tryOpen(); err != nil {
		return
	}
	return pipe.file.Read(p)
}

// Open will open the pipe if there is no file there already
func (pipe *fifoPipe) Open() (err error) {
	if pipe.opened {
		return errors.New("Pipes can't be reopened")
	}
	pipe.file, err = os.OpenFile(pipe.name, pipe.flag, os.ModeNamedPipe)
	return
}

// Name returns the name that was used to open the file as the in and out stream name
func (pipe *fifoPipe) Name() StreamNames {
	return StreamNames{
		In:  pipe.name,
		Out: pipe.name,
	}
}

// Close will close all file connections and delete all the pipes
// an attempt is made for every opperation even if the last fails
func (pipe *fifoPipe) Close() (err error) {
	err = pipe.file.Close()
	if e := os.Remove(pipe.name); e != nil {
		err = e
	}
	return
}
