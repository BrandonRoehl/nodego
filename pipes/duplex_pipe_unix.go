// +build linux darwin

package pipes

// getRandomPipeNames on unix returns two different files in temp space ending in dot pipe
func getRandomPipeNames() (pipeIn, pipeOut string, err error) {
	pipes, err := getRandomTempFiles([]string{"*.in.pipe", "*.out.pipe"})
	return pipes[0], pipes[1], err
}
