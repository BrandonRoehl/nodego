package nodego

import (
	"testing"
)

func TestPipes(t *testing.T) {

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

}
