package sorter

import (
	"bufio"
	"os"

	"tig.chaosgroup.com/candidates/bozhin.katsarski/util"
)

func NewSorter(filename string) *Sorter {
	sorter := Sorter{
		filename:  filename,
		lines:     0,
		firstByte: make([]int64, 0),
	}

	sorter.indexLines()

	return &sorter
}

type Sorter struct {
	filename  string
	lines     int     // number of lines in the file
	firstByte []int64 // index of the first byte of each line
}

func (s *Sorter) indexLines() {
	file, err := os.Open(s.filename)
	util.Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lastReadPosition := int64(0)

	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		lastReadPosition += int64(advance)
		return advance, token, err
	}
	scanner.Split(scanLines)

	currentLineStart := lastReadPosition
	for scanner.Scan() {
		s.firstByte = append(s.firstByte, currentLineStart)
		currentLineStart = lastReadPosition
	}
	s.lines = len(s.firstByte)
}

func (s *Sorter) lastByte(line int) int64 {
	if line == s.lines-1 {
		return s.fileSize()
	}
	return s.firstByte[line+1] - 1
}

func (s *Sorter) fileSize() int64 {
	fileInfo, err := os.Stat(s.filename)
	util.Check(err)
	return fileInfo.Size()
}
