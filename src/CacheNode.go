package src

import (
	"container/list"
	"time"
)

// An CacheNode is something we manage in a priority queue.
type CacheNode struct {
	// The index is needed by update and is maintained by the heap.Interface methods.
	index    int // The index of the item in the heap.
	key      string
	value    string
	expiry   time.Time
	element  *list.Element
	priority int
}
