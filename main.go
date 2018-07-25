package main

// util.CreateRandomFile("big.txt", 3200000, 500)

import (
	"fmt"
	"os"

	"tig.chaosgroup.com/candidates/bozhin.katsarski/sorter"
)

func main() {
	if len(os.Args) != 2 {
		panic("Usage: file-sorter <text-file>")
	}

	filename := os.Args[1]
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		panic(fmt.Errorf("File '%s' not found", filename))
	}

	// data, err := ioutil.ReadFile(filename)
	// util.Check(err)
	// fmt.Println(len(data))

	sorter.NewSorter(filename).Sort()

	// contents, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Print(string(contents))
}
