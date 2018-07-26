package main

import (
	"flag"
	"io/ioutil"
	"log"

	"tig.chaosgroup.com/candidates/bozhin.katsarski/gen"

	"tig.chaosgroup.com/candidates/bozhin.katsarski/sorter"
)

func main() {
	verbose, filename, fileSizeMB := getFlags()

	if !verbose {
		log.SetOutput(ioutil.Discard)
	}

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

func getFlags() (verbose bool, filename string, fileSizeMB int) {
	flag.BoolVar(&verbose, "verbose", false, "Print verbose logs")

	flag.StringVar(&filename, "filename", "", "Existing file to be sorted")

	flag.IntVar(&fileSizeMB, "filesize", 0, "Size (MB) of file to be randomly generated and sorted")

	flag.Parse()

	if (filename == "" && fileSizeMB == 0) || (filename != "" && fileSizeMB != 0) {
		panic(`Usage:
		file-sorter [-verbose] -filename=<name> # sort the existing file <name>
		file sorter [-verbose] -filesize=<size> # generate a random file of size <size> MB and sort it`)
	}
	return
}
