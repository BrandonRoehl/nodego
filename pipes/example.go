package pipes

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

func ExampleNewFifoPipe() {
	pipe, err := NewFifoPipe("test.pipe", os.O_RDWR)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Fprintln(pipe, "Hello pipe!")
		if err := pipe.Close(); err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		result, err := ioutil.ReadAll(pipe)
		if err != nil {
			panic(err)
		}
		fmt.Print(result)
	}()

	return

	// Output:
	// Hello pipe!
}
