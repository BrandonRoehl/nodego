// +build linux darwin

package pipes

import "os"

// getRandomPipeNames on unix returns two different files in temp space ending in dot pipe
func getRandomPipeNames() (pipeIn, pipeOut string, err error) {
	pipes, err := getRandomTempFiles([]string{"*.in.pipe", "*.out.pipe"})
	return pipes[0], pipes[1], err
}

// NewDuplexPipe returns an io.ReaadWriteCloser that maintains an in and an out
// pipe to write and read from for inter-process comunication
func NewDuplexPipe() (Pipe, error) {
	pipeIn, pipeOut, err := getRandomPipeNames()
	if err != nil {
		return nil, err
	}

	// READ PIPE
	inPipe, err := NewFifoPipe(pipeIn, os.O_RDONLY)
	if err != nil {
		return nil, err
	}

	// WRITE PIPE
	outPipe, err := NewFifoPipe(pipeOut, os.O_WRONLY|os.O_APPEND)
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
