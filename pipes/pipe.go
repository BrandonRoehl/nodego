package pipes

import "io"

// Pipe is the interface to interact with the pipes that are created
type Pipe interface {
	io.ReadWriteCloser

	// Name of the pipe that was created
	Name() StreamNames
}

// StreamNames contains the io named streams as file names
type StreamNames struct {
	In, Out string
}
