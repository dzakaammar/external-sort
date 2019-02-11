package exsort

import (
	"bufio"
	"fmt"
	"sort"
	"strconv"
)

//ExternalSort :nodoc:
type ExternalSort struct {
	in           string
	out          string
	partitionDir string
	size         int
}

// New :nodoc:
func New(in, out, partitionDir string, size int) *ExternalSort {
	return &ExternalSort{
		in:           in,
		out:          out,
		partitionDir: partitionDir,
		size:         size,
	}
}

//Sort :nodoc:
func (e *ExternalSort) Sort() error {
	partitionFiles, err := e.createPartition()
	if err != nil {
		return err
	}

	err = e.mergeFile(partitionFiles)
	if err != nil {
		return err
	}

	return nil
}

func (e *ExternalSort) mergeFile(partitionFiles []string) error {
	outFile, err := openFile(e.out)
	if err != nil {
		return err
	}
	defer outFile.Close()

	numberOfFiles := len(partitionFiles)
	harr := make([]*Node, numberOfFiles)
	scan := make([]*bufio.Scanner, numberOfFiles)

	for i, file := range partitionFiles {
		v, err := openFile(file)
		if err != nil {
			return err
		}
		defer v.Close()

		scanner := bufio.NewScanner(v)

		// TODO
		// Add customization
		scanner.Split(bufio.ScanLines)
		scan[i] = scanner

		if !scanner.Scan() {
			break
		}
		node := &Node{}
		node.Element, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return err
		}
		node.Index = i
		harr[i] = node
	}
	heapSize := numberOfFiles - 1
	heap := NewHeapSort(&harr, heapSize)

	count := 0
	for count != heapSize {
		root := heap.GetMin()
		err := writeToFile(outFile, root.Element)
		if err != nil {
			return err
		}

		if !scan[root.Index].Scan() {
			root.Element = int(^uint(0) >> 1)
			count++
		} else {
			root.Element, _ = strconv.Atoi(scan[root.Index].Text())
		}

		heap.ReplaceMin(root)
	}

	return outFile.Sync()
}

func (e *ExternalSort) createPartition() ([]string, error) {
	inFile, err := openFile(e.in)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	//TODO
	// Add customization on this part
	scanner.Split(bufio.ScanLines)

	arr := make([]int, e.size)

	var (
		i, numberOfFile = 0, 0
		partitionFile   []string
	)
	for ; ; numberOfFile++ {
		fn := fmt.Sprintf("%s/partition-%d", e.partitionDir, numberOfFile)
		partitionFile = append(partitionFile, fn)
		file, err := openFile(fn)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		cont := true
		for i = 0; i < e.size; i++ {
			if !scanner.Scan() {
				cont = false
				break
			}

			// TODO
			// Add customization on this part
			res, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return nil, err
			}
			arr[i] = res
		}

		// TODO
		// Add more sort choices
		sort.Ints(arr)

		for j := 0; j < i; j++ {
			err := writeToFile(file, arr[j])
			if err != nil {
				return nil, err
			}
		}
		file.Sync()

		if !cont {
			break
		}
	}

	return partitionFile, nil
}
