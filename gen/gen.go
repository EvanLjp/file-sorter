package gen

import (
	"fmt"
	"math/rand"
	"os"
)

const lineLengthBytes = 200

// RandomTextFile creates a random file with given name, size in MB, containing latin letters
func RandomTextFile(name string, sizeMB int) error {
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("RandomTextFile: %s", err)
	}
	defer f.Close()

	lines := sizeMB * 1e6 / lineLengthBytes
	for i := 0; i < lines; i++ {
		if _, err := f.WriteString(RandomString(lineLengthBytes-1) + "\n"); err != nil {
			return fmt.Errorf("RandomTextFile: %s", err)
		}
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("RandomTextFile: %s", err)
	}

	return nil
}

// RandomString generates a random string of latin letters with given length
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(97 + rand.Intn(122-97)) // a=97 and z=122
	}
	return string(bytes)
}
