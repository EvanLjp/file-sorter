package sorter

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// AssertSorted checks if file's lines are sorted.
// Panics if not so
func (s *Sorter) AssertSorted() {
	isSorted, err := isSorted(s.filename)
	if err != nil {
		panic(err)
	}

	if !isSorted {
		panic(fmt.Errorf("File '%s' is not sorted", s.filename))
	}

	log.Printf("File '%s' sorted successfully", s.filename)
}

// isSorted checks if the given file's lines are sorted
func isSorted(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("isSorted '%s': %s", filename, err)
	}

	previousLine := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() < previousLine {
			return false, nil
		}
		previousLine = scanner.Text()
	}

	if scanner.Err() != nil {
		return false, fmt.Errorf("isSorted '%s': %s", filename, err)
	}

	return true, nil
}
