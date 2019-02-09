package exsort

type Node struct {
	Element int
	Index   int
}

type Heap struct {
	Nodes []*Node
	Size  int
}

func NewHeapSort(nodes *[]*Node, size int) *Heap {
	heap := &Heap{
		Nodes: *nodes,
		Size:  size,
	}

	for i := (heap.Size - 1) / 2; i >= 0; i-- {
		heap.Heapify(i)
	}

	return heap
}

func (heap *Heap) RemoveTop(length int) {
	var lastIndex = length - 1
	heap.Nodes[0], heap.Nodes[lastIndex] = heap.Nodes[lastIndex], heap.Nodes[0]
	heap.Heapify(0)
}

func (heap *Heap) Heapify(root int) {
	smallest := root
	l, r := heap.Left(root), heap.Right(root)

	if l < heap.Size && heap.Nodes[l].Element < heap.Nodes[root].Element {
		smallest = l
	}

	if r < heap.Size && heap.Nodes[r].Element < heap.Nodes[smallest].Element {
		smallest = r
	}

	if smallest != root {
		heap.Nodes[root], heap.Nodes[smallest] = heap.Nodes[smallest], heap.Nodes[root]
		heap.Heapify(smallest)
	}
}

func (*Heap) Left(root int) int {
	return (root * 2) + 1
}

func (*Heap) Right(root int) int {
	return (root * 2) + 2
}

func (h *Heap) GetMin() *Node {
	return h.Nodes[0]
}

func (h *Heap) ReplaceMin(x *Node) {
	h.Nodes[0] = x
	h.Heapify(0)
}
