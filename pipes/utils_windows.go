//build+ windows

package pipes

func tempDir() string {
	return `\\.\pipe\`
}
