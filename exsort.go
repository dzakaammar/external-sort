package exsort

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

//ExternalSort :nodoc:
type ExternalSort struct {
	inputFile  *os.File
	outputFile *os.File
	opts       Options
}

//SortOptions :nodoc:
type SortOptions func(*Options) error

//Options :nodoc:
type Options struct {
	Size int
}

//DefaultOptions :nodoc:
var DefaultOptions = Options{
	Size: 100,
}

//SetPartitionSize :nodoc:
func SetPartitionSize(size int) SortOptions {
	return func(o *Options) error {
		o.Size = size
		return nil
	}
}

//Sort :nodoc:
func (e *ExternalSort) Sort() error {
	numOfPartition, err := e.createPartition()
	if err != nil {
		return err
	}

	err = e.mergeFile(numOfPartition)
	if err != nil {
		return err
	}
	return nil
}

func (e *ExternalSort) mergeFile(numOfPartition int) error {
	harr := make([]*Node, numOfPartition)
	scan := make([]*bufio.Scanner, numOfPartition)

	for i := 0; i < numOfPartition; i++ {
		v, err := e.openFile(fmt.Sprintf("./partition-%d", i))
		if err != nil {
			return err
		}
		defer v.Close()

		scanner := bufio.NewScanner(v)
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
	heapSize := numOfPartition - 1
	heap := NewHeapSort(&harr, heapSize)

	count := 0
	for count != heapSize {
		root := heap.GetMin()
		err := e.writeToFile(e.outputFile, root.Element)
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

	return e.outputFile.Sync()
}

func (e *ExternalSort) openFile(path string) (*os.File, error) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (e *ExternalSort) writeToFile(file *os.File, data int) error {
	_, err := file.WriteString(fmt.Sprintf("%d\n", data))
	if err != nil {
		return err
	}

	return nil
}

func (e *ExternalSort) createPartition() (int, error) {
	numberOfFile := 0

	scanner := bufio.NewScanner(e.inputFile)
	scanner.Split(bufio.ScanLines)
	arr := make([]int, e.opts.Size)

	var i int
	for ; ; numberOfFile++ {
		file, err := e.openFile(fmt.Sprintf("./partition-%d", numberOfFile))
		if err != nil {
			return 0, err
		}
		defer file.Close()

		cont := true
		for i = 0; i < e.opts.Size; i++ {
			if !scanner.Scan() {
				cont = false
				break
			}

			res, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return 0, err
			}
			arr[i] = res
		}

		sort.Ints(arr)

		for j := 0; j < i; j++ {
			err := e.writeToFile(file, arr[j])
			if err != nil {
				return 0, err
			}
		}
		file.Sync()

		if !cont {
			break
		}
	}

	return numberOfFile, nil
}
