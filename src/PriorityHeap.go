package src

// This code is taken from the https://pkg.go.dev/container/heap int heap example

// An PriorityHeap is a min-heap of ints.
type PriorityHeap []int

func (h PriorityHeap) Len() int           { return len(h) }
func (h PriorityHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h PriorityHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PriorityHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *PriorityHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
