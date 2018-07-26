package sorter

import (
	"bufio"
	"fmt"
	"os"

	"tig.chaosgroup.com/candidates/bozhin.katsarski/util"
)

// Sorter encapsulates the data needed for sorting a file
type Sorter struct {
	filename       string
	fileSizeBytes  int64
	lines          int     // number of lines in the file
	linesFirstByte []int64 // index of the first byte of each line
}

// New creates a new file sorter
func New(filename string) (*Sorter, error) {
	if !util.FileExists(filename) {
		return nil, fmt.Errorf("Could not create sorter, file '%s' was not found", filename)
	}

	sorter := Sorter{
		filename:       filename,
		lines:          0,
		linesFirstByte: make([]int64, 0),
	}

	fileSizeBytes, err := util.FileSize(filename)
	if err != nil {
		panic(err)
	}
	sorter.fileSizeBytes = fileSizeBytes

	sorter.indexLines()

	return &sorter, nil
}

// indexLines reads the file line by line and indexes the first byte of each line
func (s *Sorter) indexLines() error {
	file, err := os.Open(s.filename)
	if err != nil {
		return fmt.Errorf("indexLines: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// we use byteIndexToBeRead from scanLines's closure to know where each lines begins
	byteIndexToBeRead := int64(0)
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		byteIndexToBeRead += int64(advance)
		return advance, token, err
	}

	scanner.Split(scanLines)

	currentLineFirstByte := int64(0)
	for scanner.Scan() {
		s.linesFirstByte = append(s.linesFirstByte, currentLineFirstByte)

		currentLineFirstByte = byteIndexToBeRead
	}

	if scanner.Err() != nil {
		return fmt.Errorf("indexLines: %s", scanner.Err())
	}

	s.lines = len(s.linesFirstByte)
	return nil
}

// firstByte returns the index of the first byte of the given line
func (s *Sorter) firstByte(lineIndex int) int64 {
	return s.linesFirstByte[lineIndex]
}

// lastByte returns the index of the last byte of the given line
func (s *Sorter) lastByte(lineIndex int) int64 {
	if lineIndex == s.lines-1 {
		return s.fileSizeBytes
	}
	return s.linesFirstByte[lineIndex+1] - 1
}
