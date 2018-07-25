package util

import (
	"encoding/base64"
	"math/rand"
	"os"
)

func CreateRandomFile(filename string, lineCount int, lineLengthBytes int) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	line := make([]byte, lineLengthBytes, lineLengthBytes)

	for i := 0; i < lineCount; i++ {
		rand.Read(line)

		_, err := f.WriteString(base64.StdEncoding.EncodeToString(line) + "\n")
		if err != nil {
			panic(err)
		}
	}
	f.Sync()
}
