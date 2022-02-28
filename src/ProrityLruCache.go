package src

import (
	"container/heap"
	"container/list"
	"time"
)

type PriorityLruCache struct {
	keyValuePair   map[string]*CacheNode
	maxSize        int
	timestampQueue TimestampQueue
	priorityMap    map[int]*list.List
	priorityHeap   PriorityHeap
}

// New Creates a new PriorityLruCache
func New(maxSize int) PriorityLruCache {
	timestampQueue := make(TimestampQueue, 0)
	heap.Init(&timestampQueue)
	priorityHeap := PriorityHeap{}
	heap.Init(&priorityHeap)
	return PriorityLruCache{
		keyValuePair:   make(map[string]*CacheNode),
		maxSize:        maxSize,
		timestampQueue: timestampQueue,
		priorityMap:    make(map[int]*list.List),
		priorityHeap:   priorityHeap,
	}

}

// Set will do the following in this order
// 1. If the element already exists, update the value and move it to the front of the LRU cache
//		1a. if we are updating the priority, update the map for that as well
//		1b. if we are updating the expiry, update the queue for that as well
// 2. If the queue is full
//		2a. Check if the oldest element and see if its expired and if it is, remove it.
//		2b. if 2a did not change the length, find the lowest priority element.
//			2bb. If there are multiple elements remove the least recently used element.
// 3. Add the new element
func (lru *PriorityLruCache) Set(key string, value string, priority int, expiry time.Time) bool {
	if node, ok := lru.keyValuePair[key]; ok {
		// update value
		node.value = value

		// if we update the priority remove it from its all priority map
		if node.priority != priority {
			lru.priorityMap[node.priority].Remove(node.element)
			node.priority = priority
		}

		// if we have updated the expiry update it in the timestampQueue as well
		if node.expiry != expiry {
			lru.timestampQueue.update(node, value, expiry)
			node.expiry = expiry
		}

		// put to front of its priority queue
		lru.priorityMap[node.priority].MoveToFront(node.element)
	}
	// the cache is full evict one item
	if lru.isFull() {
		if lru.maxSize == 0 {
			return false
		}
		lru.evict()
		// if the cache is still full, remove the lowest priority element.
		if lru.isFull() {

			lowestPriority := lru.priorityHeap.Pop().(int)
			leastRecentlyUsedValue := lru.priorityMap[lowestPriority].Back()
			node := leastRecentlyUsedValue.Value.(*CacheNode)
			lru.remove(node)
			lru.timestampQueue.remove(node)
		}
	}
	// add the new item
	lru.add(key, value, priority, expiry)
	return true
}

// Get will return True if the element exists and the value. This will also move the element
// to the most recently used
func (lru PriorityLruCache) Get(key string) (string, bool) {
	// find the value
	if node, ok := lru.keyValuePair[key]; ok {
		// if its expired return False
		if node.expiry.Before(time.Now()) {
			return "", false
		}
		lru.priorityMap[node.priority].MoveToFront(node.element)

		return node.value, true
	}
	return "", false
}

// evict Checks the oldest element and removes it if its expired
func (lru *PriorityLruCache) evict() {
	// pop value from timestampHeap.
	node := heap.Pop(&lru.timestampQueue).(*CacheNode)
	// If the value is not expired put it back
	if node.expiry.After(time.Now()) {
		heap.Push(&lru.timestampQueue, node)
	} else {
		lru.remove(node)
	}

}

// isFull returns true if the PriorityLruCache is full
func (lru PriorityLruCache) isFull() bool {
	return len(lru.keyValuePair) >= lru.maxSize
}

// remove Removes an element from priorityMap and keyValuePair but not timestampQueue that is done
// in set and evict for efficient reasons
func (lru PriorityLruCache) remove(node *CacheNode) bool {
	if value, ok := lru.keyValuePair[node.key]; ok {
		// remove it from the priorityMap
		lru.priorityMap[value.priority].Remove(node.element)
		if lru.priorityMap[value.priority].Len() == 0 {
			delete(lru.priorityMap, value.priority)
		}
		// remove it from the keyValuePair Map
		delete(lru.keyValuePair, value.key)
	}
	return false

}

// add Adds the value to the PriorityLruCache. Unless its full, then it will return False
func (lru *PriorityLruCache) add(key string, value string, priority int, expiry time.Time) bool {
	if lru.isFull() {
		return false
	}

	node := &CacheNode{
		key:      key,
		value:    value,
		expiry:   expiry,
		priority: priority,
	}

	// add the value to timeStamp Queue
	heap.Push(&lru.timestampQueue, node)

	// add the value to the priority map
	if _, ok := lru.priorityMap[priority]; !ok {
		heap.Push(&lru.priorityHeap, priority)
		lru.priorityMap[priority] = list.New()
	}

	element := lru.priorityMap[priority].PushFront(node)
	node.element = element

	// add the value to the value map
	lru.keyValuePair[key] = node

	return true
}
