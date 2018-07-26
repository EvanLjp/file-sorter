package sorter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

// Sort sorts the lines of the file by overwriting it
func (s *Sorter) Sort() error {
	maxBytesPerBatch := calcMaxBytesPerBatch(s.fileSizeBytes)

	batchesCount, err := s.sortByBatches(maxBytesPerBatch)
	if err != nil {
		return fmt.Errorf("Sort: %s", err)
	}

	defer deleteTmpFiles(batchesCount)

	if err := s.mergeBatches(batchesCount); err != nil {
		return fmt.Errorf("Sort: %s", err)
	}

	return nil
}

// sortByBatches divides the file's lines into batches with a fixed maximum size
// and in-memory sorts each batch into a temporary file. Returns number of batches used and error
func (s *Sorter) sortByBatches(maxBytesPerBatch int64) (int, error) {
	log.Printf("Sorting with max batch size %d bytes...", maxBytesPerBatch)

	currentBatch := 0
	fromLine := 0
	for fromLine < s.lines {
		toLine := fromLine + 1
		for toLine < s.lines && s.lastByte(toLine)-s.firstByte(fromLine) <= maxBytesPerBatch {
			toLine++
		}

		if err := s.sortBatchIntoFile(fromLine, toLine, tmpFilename(currentBatch)); err != nil {
			return 0, fmt.Errorf("sortByBatches: %s", err)
		}

		currentBatch++
		fromLine = toLine
	}

	return currentBatch, nil
}

// sortBatchIntoFile sorts (in-memory) the lines with indices [fromLine; toLine)
// into a new file with the given name.
// It truncates the file if it already exists
func (s *Sorter) sortBatchIntoFile(fromLine, toLine int, outputFile string) error {
	log.Printf("Batch sorting bytes [%d, %d)...", s.firstByte(fromLine), s.lastByte(toLine-1))

	file, err := os.Open(s.filename)
	if err != nil {
		return fmt.Errorf("sortBatchIntoFile: %s", err)
	}
	defer file.Close()

	_, err = file.Seek(s.firstByte(fromLine), 0)
	if err != nil {
		return fmt.Errorf("sortBatchIntoFile lines [%d, %d)", fromLine, toLine)
	}

	lines := make([]string, 0, toLine-fromLine)
	scanner := bufio.NewScanner(file)
	for i := fromLine; i < toLine; i++ {
		ok := scanner.Scan()
		if !ok {
			return fmt.Errorf("sortBatchIntoFile: Could not read expected number of lines from file: %s", scanner.Err())
		}

		lines = append(lines, scanner.Text())
	}

	sort.Strings(lines)

	outFile, err := os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("sortBatchIntoFile: %s", err)
	}
	defer outFile.Close()

	for i := range lines {
		_, err := outFile.WriteString(lines[i] + "\n")
		if err != nil {
			return fmt.Errorf("sortBatchIntoFile: %s", err)
		}
	}

	return nil
}

// mergeBatches merges the batches from tmp files into the original file
func (s *Sorter) mergeBatches(batchesCount int) error {
	log.Printf("Merging %d batches...", batchesCount)

	scanners := make([]*bufio.Scanner, 0, batchesCount)

	for i := 0; i < batchesCount; i++ {
		file, err := os.Open(tmpFilename(i))
		if err != nil {
			return fmt.Errorf("mergeBatches: %s", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		scanners = append(scanners, scanner)
	}

	file, err := os.OpenFile(s.filename, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("mergeBatches: %s", err)
	}

	for {
		minLineIndex := -1 // index of the scanner with the smallest line
		minLine := ""
		for i := 0; i < batchesCount; i++ {
			if scanners[i] != nil {
				if minLineIndex == -1 || scanners[i].Text() < minLine {
					minLineIndex = i
					minLine = scanners[i].Text()
				}
			}
		}

		if minLineIndex == -1 {
			// all scanners have finished
			break
		}

		file.WriteString(minLine + "\n")

		if !scanners[minLineIndex].Scan() {
			scanners[minLineIndex] = nil
		}
	}

	return nil
}

// tmpFilename returns the name of the tmp file with given index
func tmpFilename(tmpFileIndex int) string {
	const tmpFileSuffix = ".file-sorter.tmp"
	return fmt.Sprintf("%d%s", tmpFileIndex, tmpFileSuffix)
}

// deleteTmpFiles deletes the tmp files created during sorting
func deleteTmpFiles(batchesCount int) {
	for i := 0; i < batchesCount; i++ {
		err := os.Remove(tmpFilename(i))
		if err != nil {
			fmt.Printf("deleteTmpFiles: file '%s': %s", tmpFilename(i), err)
		}
	}
}

func calcMaxBytesPerBatch(fileSize int64) int64 {
	return 1e3 // TODO
}
