# file-sorter

A cli tool that sorts files too big to fit in RAM.

## Algorithm

File sorter uses an [external sort  algorithm](https://en.wikipedia.org/wiki/External_sorting) (similar to merge sort) in order not to load the whole input file into memory.

* First the file is read line by line and an index is created in order to know at which byte each line begins.
* Then, the file is divided into _batches_ of fixed maximal size, each of which containing a few consequtive lines of the file.
* One by one, each such batch is loaded into memory and sorted using a conventional in-memory sorting algorithm (for example quick sort) and then written to a temporary file on the disk.
* Note batches are processed one by one in order not to occupy too much memory. Computation time is traded off for less memory consumption.
* After all __K__ batches are sorted in temporary files, they are all merged into the original file, using a [__K__-way merge algorithm](https://en.wikipedia.org/wiki/K-way_merge_algorithm).
* This means each of the temporary files starts to be read line by line as if it were a queue, then at each step of the merging process, the lines at the front of each queue are considered and the least one is chosen to be written to the original file. This continues until all lines are merged.
* In this implementation of the algorithm at each step of the merging process the correct line is chosen using a linear search over the fronts of the mentioned queues, which is suboptimal.
* The algorithm can be improved if instead pointers to the queues' fronts are kept in a heap (priority queue). That would achieve O(__N__ * log __K__) complexity instead of O(__N__ * __K__) where __N__ is the total number of lines in the original file and __K__ - the number of batches used. However, in most cases __K__ is not big enough to justify the added complexity.
* After it is finished, the algorithm checks that the file is sorted correctly

## Usage

Build the tool: `go build -o file-sorter`.

See its usage: `./file-sorter`:
* `file-sorter [-batch=<maxBatchSizeB>] -filename=<name>` sort an existing file
* `file-sorter [-batch=<maxBatchSizeB>] -filesize=<sizeMB>` generate a random file with given size in MB and then sort it

The user can specify the max size (bytes) of batches to use while sorting (1e8 = 100MB by default)

Note that generating a big random file is quite slow. The user may want to generate it once and use it multiple times for testing.

Also if the input file is too many times (>1000) bigger than the batch size,
too many tmp files for the batches will be created, which may be a problem for the os.
Also such a ratio would not be optimal.

## Benchmarks

On a mid-range laptop, the following results were obtained:

| File size (MB) | Max batch size (MB) | Approx time (Sec) |
|----------------|---------------------|-------------------|
| 100            | 10                  | 56                |
| 100            | 100                 | 4                 |
|                |                     |                   |
| 1 000          | 10                  | 95                |
| 1 000          | 100                 | 50                |
| 1 000          | 1 00                | 68                |
|                |                     |                   |
| 10 000         | 1 00                | 1620              |

Note that the max batch size is not the max total memory that the file sorter needs.
For example when sorting a 1 000MB file with max batch size  100MB, it consumed about 300MB memory.
