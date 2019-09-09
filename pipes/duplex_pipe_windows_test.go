//build+ windows

package pipes

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndCloseDuplex(t *testing.T) {
	testMsg := []byte("test message")

	pipe, err := NewDuplexPipe()
	assert.NoError(t, err, "Expected error to be nil when creating pipe")

	// Defer close for a premeture test termination cleanup
	defer pipe.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// Attempt to read from the pipe
	wg.Done()
	buffer := make([]byte, len(testMsg))
	n, err := pipe.Read(buffer)

	assert.NoError(t, err, "Expected no error when reading from the pipe")
	assert.Equal(t, len(testMsg), n)
	assert.Equal(t, testMsg, buffer, "Expected the message to be what was sent")
	t.Log("Reading passed")

	// Attempt to write to the pipe
	go func() {
		defer wg.Done()
		n, err := pipe.Write(testMsg)
		assert.NoError(t, err, "Expected no error when writing to the pipe")
		assert.Equal(t, len(testMsg), n)
		t.Log("Writting passed")
	}()

	wg.Wait()

	assert.NoError(t, pipe.Close(), "Pipe can be closed")

}
