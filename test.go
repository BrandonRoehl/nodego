package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"github.com/brandonroehl/nodego/pipes"
)

func nonNilPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pipe, err := pipes.NewDuplexPipe()
	nonNilPanic(err)
	defer pipe.Close()

	fmt.Println("Out:", pipe.Name().Out)
	fmt.Println("In:", pipe.Name().In)

	// Take the pipes and turn it to a decoder
	reader := bufio.NewReader(pipe)
	dec := json.NewDecoder(reader)

	// Take the pipes and turn them into an encoder
	enc := json.NewEncoder(pipe)

	// Runtime loop to decode json objects into the interface
	for dec.More() {
		// 1. Read in the object as JSON
		var i interface{}
		if err := dec.Decode(&i); err != nil {
			log.Fatal(err)
		}

		// 2. Do something with the json
		log.Println(&i, i)

		// 3. Print JSON to the out stream for the result
		enc.Encode(i)
	}
}
