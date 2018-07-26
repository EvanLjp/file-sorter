package main

import (
	"flag"

	"tig.chaosgroup.com/candidates/bozhin.katsarski/gen"

	"tig.chaosgroup.com/candidates/bozhin.katsarski/sorter"
)

func main() {
	filename, fileSizeMB := getFlags()

	if fileSizeMB != 0 {
		filename = "random-file.txt"
		gen.RandomTextFile(filename, fileSizeMB)
	}

	s, err := sorter.New(filename)
	if err != nil {
		panic(err)
	}

	if err := s.Sort(); err != nil {
		panic(err)
	}

	s.AssertSorted()
}

func getFlags() (filename string, fileSizeMB int) {
	flag.StringVar(&filename, "filename", "", "Existing file to be sorted")

	flag.IntVar(&fileSizeMB, "filesize", 0, "Size (MB) of file to be randomly generated and sorted")

	flag.Parse()

	if (filename == "" && fileSizeMB == 0) || (filename != "" && fileSizeMB != 0) {
		panic(`Usage:
		file-sorter -filename=<name> # sort the existing file <name>
		file sorter -filesize=<size> # generate a random file of size <size> MB and sort it`)
	}
	return
}
