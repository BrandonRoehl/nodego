// +build linux darwin

package pipes

import "os"

func tempDir() string {
	return os.TempDir()
}
