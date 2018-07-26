# file-sorter

A command line tool that sorts the lines of large text files.

## Algorithm

File sorter uses an [external sort](https://en.wikipedia.org/wiki/External_sorting) algorithm in order not to load the whole input file into memory.

* First the file is read line by line and an index is created in order to know at which byte each line begins.
* Then, the file is divided into _batches_ of fixed maximal size, each of which containing a few consequtive lines of the file.
* One by one, each such batch is loaded into memory and sorted using a conventional in-memory sorting algorithm (for example quick sort) and then written to a temporary file on the disk.
* Note batches are processed one by one in order not to occupy too much memory. Computation time is traded off for less memory consumption.
* After all __K__ batches are sorted in temporary files, they are all merged into the original file, using a [__K__-way merge algorithm](https://en.wikipedia.org/wiki/K-way_merge_algorithm).
* This means each of the temporary files starts to be read line by line as if it were a queue, then at each step of the merging process, the lines at the front of each queue are considered and the least one is chosen to be written to the original file. This continues until all lines are merged.
* In this implementation of the algorithm at each step of the merging process the correct line is chosen using a linear search over the fronts of the mentioned queues, which is suboptimal.
* The algorithm can be improved if instead pointers to the queues' fronts are kept in a heap (priority queue). That would achieve O(__N__ * log __K__) complexity instead of O(__N__ * __K__) where __N__ is the total number of lines in the original file and __K__ - the number of batches used. However, in most cases __K__ is not big enough to justify the added complexity.

## Usage

Build the tool: `go build -o file-sorter`.

See its usage: `./file-sorter`:
* `file-sorter -filename=<name>` sort an existing file
* `file-sorter -filesize=<sizeMB>` generate a random file with given size in MB and then sort it

Note that generating a big random file is quite slow.
