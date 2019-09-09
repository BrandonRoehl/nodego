//build+ windows

package pipes

import (
	"errors"
	"net"
	"os"

	"github.com/microsoft/go-winio"
)

// NewDuplexPipe returns an io.ReaadWriteCloser that maintains an in and an out
// pipe to write and read from for inter-process comunication
func NewDuplexPipe() (Pipe, error) {
	pipeNames, err := getRandomTempFiles([]string{"*.pipe"})
	if err != nil {
		return nil, err
	}
	return NewNamedDuplexPipe(pipeNames[0], "")
}

// NewNamedDuplexPipe returns an io.ReaadWriteCloser that maintains an in and an out
// pipe to write and read from for inter-process comunication
func NewNamedDuplexPipe(name, _ string) (Pipe, error) {
	return &duplexPipe{
		in:     nil,
		out:    nil,
		name:   name,
		closed: false,
	}, nil
}

type duplexPipe struct {
	in, out net.Conn
	name    string
	closed  bool
}

// Write will try and connect to a pipe if this is the first time or delegate to the file writer
func (pipe *duplexPipe) Write(p []byte) (n int, err error) {
	if pipe.closed {
		err = errors.New("Pipe has already been closed")
		return
	}

	if pipe.out == nil {
		// Open the write pipe
		pipe.out, err = winio.DialPipe(pipe.name, nil)
		if err != nil {
			return
		}
	}

	return pipe.out.Write(p)
}

// Read will try and connect to a pipe if this is the first time or delegate to the file reader
func (pipe *duplexPipe) Read(p []byte) (n int, err error) {
	if pipe.closed {
		err = errors.New("Pipe has already been closed")
		return
	}

	if pipe.out == nil {
		// Open the read pipe
		var l net.Listener
		l, err = winio.ListenPipe(pipe.name, nil)
		if err != nil {
			return
		}

		pipe.in, err = l.Accept()
		if err != nil {
			return
		}
	}

	return pipe.in.Read(p)
}

// Open will open the pipe if there is no file there already
func (pipe *duplexPipe) Open() error {
	return nil
}

// Name returns the name that was used to open the file as the in and out stream name
func (pipe *duplexPipe) Name() StreamNames {
	return StreamNames{
		In:  pipe.name,
		Out: pipe.name,
	}
}

// Close will close all file connections and delete all the pipes
// an attempt is made for every opperation even if the last fails
func (pipe *duplexPipe) Close() (err error) {
	if pipe.closed {
		return errors.New("Pipe is already closed")
	}

	err = pipe.in.Close()
	if e := pipe.out.Close(); e != nil {
		err = e
	}
	if e := os.Remove(pipe.name); e != nil {
		err = e
	}
	return
}
