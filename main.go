package main

import (
	"flag"
	"log"

	"github.com/file-sorter/gen"

	"github.com/file-sorter/sorter"
)

func main() {
	filename, fileSizeMB, maxBatchSizeB := getFlags()

	if fileSizeMB != 0 {
		filename = "random.txt"
		log.Printf("Generating random %d MB file '%s'...", fileSizeMB, filename)
		gen.RandomTextFile(filename, fileSizeMB)
		log.Println("Random file generated")
	}

	s, err := sorter.New(filename)
	if err != nil {
		panic(err)
	}

	if maxBatchSizeB != 0 {
		s.MaxBatchSizeB(maxBatchSizeB)
	}

	if err := s.Sort(); err != nil {
		panic(err)
	}

	s.AssertSorted()
}

func getFlags() (filename string, fileSizeMB int, maxBatchSizeB int64) {
	flag.StringVar(&filename, "filename", "", "Existing file to be sorted")

	flag.IntVar(&fileSizeMB, "filesize", 0, "Size (MB) of file to be randomly generated and sorted")

	flag.Int64Var(&maxBatchSizeB, "batch", 0, "Max size (bytes) of batches to sort in")

	flag.Parse()

	if (filename == "" && fileSizeMB == 0) || (filename != "" && fileSizeMB != 0) {
		panic(`Usage:
		file-sorter [-batch=<maxBatchSizeB>] -filename=<name> # sort the existing file <name>
		file sorter [-batch=<maxBatchSizeB>] -filesize=<size> # generate a random file of size <size> MB and sort it`)
	}
	return
}
