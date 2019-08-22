package nodego

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fmt.Println(pipeIn)
// fmt.Println()

// if err := unix.Mkfifo(pipeIn, 0600); err != nil {
// 	panic(err)
// }

// // READ PIPE
// file, err := os.OpenFile(pipeIn, os.O_CREATE, os.ModeNamedPipe)
// if err != nil {
// 	panic(err)
// }
// // Cleanup the pipe whenever it is closed
// defer func() {
// 	file.Close()
// 	if err := os.Remove(pipeIn); err != nil {
// 		panic(err)
// 	}
// }()

// // Read from the pipe until the pipe is closed
// reader := bufio.NewReader(file)
// dec := json.NewDecoder(reader)

// // Runtime loop to decode json objects into the interface
// for dec.More() {
// 	var i interface{}
// 	if err := dec.Decode(&i); err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println(&i, i)
// }

func TestSinglePipe(t *testing.T) {
	testMsg := []byte("test message")

	pipeName, _, err := getRandomPipeNames()
	assert.NoError(t, err, "Expected error to be nil when getting names")

	pipe, err := makePipe(pipeName, os.O_RDWR|os.O_CREATE)
	assert.NoError(t, err, "Expected error to be nil when creating pipe")

	// Defer close for a premeture test termination cleanup
	defer pipe.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// Attempt to write to the pipe
	go func() {
		defer wg.Done()
		n, err := pipe.Write(testMsg)
		assert.NoError(t, err, "Expected no error when writing to the pipe")
		assert.Equal(t, len(testMsg), n)
		t.Log("Writting passed")
	}()

	// Attempt to read from the pipe
	go func() {
		defer wg.Done()
		buffer := make([]byte, len(testMsg))
		n, err := pipe.Read(buffer)

		assert.NoError(t, err, "Expected no error when reading from the pipe")
		assert.Equal(t, len(testMsg), n)
		assert.Equal(t, testMsg, buffer, "Expected the message to be what was sent")
		t.Log("Reading passed")
	}()

	wg.Wait()

	assert.NoError(t, pipe.Close(), "Pipe can be closed")
}
