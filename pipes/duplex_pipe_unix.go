// +build linux darwin

package pipes

import "os"

// NewDuplexPipe returns an io.ReaadWriteCloser that maintains an in and an out
// pipe to write and read from for inter-process comunication
func NewDuplexPipe() (Pipe, error) {
	pipeIn, pipeOut, err := getRandomPipeNames()
	if err != nil {
		return nil, err
	}
	return NewNamedDuplexPipe(pipeIn, pipeOut)
}

// NewNamedDuplexPipe returns an io.ReaadWriteCloser that maintains an in and an out
// pipe to write and read from for inter-process comunication
func NewNamedDuplexPipe(inFile, outFile string) (Pipe, error) {
	// READ PIPE
	inPipe, err := NewFifoPipe(inFile, os.O_RDONLY)
	if err != nil {
		return nil, err
	}

	// WRITE PIPE
	outPipe, err := NewFifoPipe(outFile, os.O_WRONLY|os.O_APPEND)
	if err != nil {
		return nil, err
	}

	return &duplexPipe{
		inPipe:  inPipe,
		outPipe: outPipe,
	}, nil
}

type duplexPipe struct {
	inPipe, outPipe Pipe
}

func (pipe *duplexPipe) Name() StreamNames {
	return StreamNames{
		In:  pipe.inPipe.Name().In,
		Out: pipe.outPipe.Name().Out,
	}
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
