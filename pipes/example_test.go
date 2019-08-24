package pipes

import (
	"bufio"
	"fmt"
	"os"
)

func ExampleNewFifoPipe() {
	// Create a new pipe
	pipe, err := NewFifoPipe("test.pipe", os.O_RDWR)
	if err != nil {
		panic(err)
	}

	// Process 1 - Prints to the pipe
	go func() {
		fmt.Fprintln(pipe, "Hello pipe!")
	}()

	// Process 2 - Reads from the pipe
	reader := bufio.NewReader(pipe)
	result, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Closing the pipe will also delete the pipe of the same name created
	if err := pipe.Close(); err != nil {
		panic(err)
	}

	// Output:
	// Hello pipe!
}
