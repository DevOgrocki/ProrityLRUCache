package src

import (
	"container/heap"
	"time"
)

// This code is taken from the https://pkg.go.dev/container/heap PriorityQueue example

// A TimestampQueue implements heap.Interface and holds Items.
type TimestampQueue []*CacheNode

func (pq TimestampQueue) Len() int { return len(pq) }

func (pq TimestampQueue) Less(i, j int) bool {
	if pq[i] != nil && pq[j] != nil {
		return pq[i].expiry.Before(pq[j].expiry)
	}
	return false
}

func (pq TimestampQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *TimestampQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*CacheNode)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *TimestampQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *TimestampQueue) update(item *CacheNode, value string, expiry time.Time) {
	item.value = value
	item.expiry = expiry
	heap.Fix(pq, item.index)
}

func (pq *TimestampQueue) remove(node *CacheNode) {
	pq.update(node, node.value, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	heap.Pop(pq)
}
