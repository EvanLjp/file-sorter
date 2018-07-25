package sorter

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"tig.chaosgroup.com/candidates/bozhin.katsarski/util"
)

const maxBytesPerBatch = 1e3

func (s *Sorter) Sort() {
	batch := 0
	fromLine := 0
	for fromLine < s.lines {
		toLine := fromLine + 1
		for toLine < s.lines && s.lastByte(toLine)-s.firstByte[fromLine] <= maxBytesPerBatch {
			toLine++
		}

		s.sortBatch(fromLine, toLine, fmt.Sprintf("batch-%d.file-sorter.tmp", batch))
		batch++
		fromLine = toLine
	}

	fmt.Println(s.fileSize())
}

// sortBatch sorts the lines with indices [fromLine; toLine) into a new file with given name
func (s *Sorter) sortBatch(fromLine, toLine int, output string) {
	fmt.Printf("Batch sorting bytes [%d, %d)\n", s.firstByte[fromLine], s.lastByte(toLine-1))

	file, err := os.Open(s.filename)
	defer file.Close()
	util.Check(err)
	file.Seek(s.firstByte[fromLine], 0)

	lines := make([]string, 0, toLine-fromLine)
	scanner := bufio.NewScanner(file)
	for i := fromLine; i < toLine; i++ {
		ok := scanner.Scan()
		if !ok {
			fmt.Println(scanner.Err())
		} else {
			lines = append(lines, scanner.Text())
		}
	}

	sort.Strings(lines)

	outFile, err := os.OpenFile(output, os.O_EXCL|os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	defer outFile.Close()
	util.Check(err)
	for i := range lines {
		_, err := outFile.WriteString(lines[i] + "\n")
		util.Check(err)
	}
}
