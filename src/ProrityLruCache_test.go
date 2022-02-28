package src

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPriorityLruCache_New(t *testing.T) {
	lruCache := New(5)
	assert.Equal(t, 5, lruCache.maxSize)
}

func TestPriorityLruCache_Get(t *testing.T) {
	lruCache := New(5)

	lruCache.Set("TestKey", "TestValue", 1, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	actual, exists := lruCache.Get("TestKey")
	assert.Equal(t, "TestValue", actual)
	assert.True(t, exists)
}

func TestPriorityLruCache_GetMovesValueToFront(t *testing.T) {
	lruCache := New(5)

	lruCache.Set("SecondKey", "SecondValue", 1, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	lruCache.Set("FirstKey", "FirstValue", 1, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))

	assert.Equal(t, "FirstValue", lruCache.priorityMap[1].Front().Value.(*CacheNode).value)

	lruCache.Get("SecondKey")

	assert.Equal(t, "SecondValue", lruCache.priorityMap[1].Front().Value.(*CacheNode).value)
}

func TestPriorityLruCache_Set(t *testing.T) {
	lruCache := New(5)
	lruCache.Set("TestKey", "TestValue", 1, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))

	assert.Equal(t, "TestValue", lruCache.keyValuePair["TestKey"].value)
	assert.Equal(t, "TestValue", lruCache.timestampQueue[0].value)
	assert.Equal(t, "TestValue", lruCache.priorityMap[1].Front().Value.(*CacheNode).value)
}

func TestPriorityLruCache_SetRemoveLowestPriority(t *testing.T) {
	lruCache := New(3)
	lruCache.Set("TestKey1", "TestValue1", 1, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	lruCache.Set("TestKey2", "TestValue2", 2, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	lruCache.Set("TestKey3", "TestValue3", 3, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	lruCache.Set("TestKey4", "TestValue4", 4, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))

	assert.Equal(t, "TestValue4", lruCache.keyValuePair["TestKey4"].value)
	assert.Equal(t, "TestValue4", lruCache.timestampQueue[2].value)
	assert.Equal(t, "TestValue4", lruCache.priorityMap[4].Front().Value.(*CacheNode).value)
}

func TestPriorityLruCache_SetRemoveLeastRecentlyUsed(t *testing.T) {
	lruCache := New(3)
	// ascii example []
	lruCache.Set("TestKey1", "TestValue1", 1, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	// [1]
	lruCache.Set("TestKey2", "TestValue2", 1, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	// [2,1]
	lruCache.Set("TestKey3", "TestValue3", 1, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	// [3,2,1]
	lruCache.Get("TestKey1")
	// [1,3,2]
	lruCache.Set("TestKey4", "TestValue4", 1, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	// [4,1,3]

	assert.Equal(t, "TestValue4", lruCache.keyValuePair["TestKey4"].value)
	assert.Equal(t, "TestValue4", lruCache.timestampQueue[2].value)

	firstElement := lruCache.priorityMap[1].Front()
	lruCache.priorityMap[1].Remove(firstElement)
	secondElement := lruCache.priorityMap[1].Front()
	lruCache.priorityMap[1].Remove(secondElement)
	thirdElement := lruCache.priorityMap[1].Front()
	lruCache.priorityMap[1].Remove(thirdElement)

	assert.Equal(t, "TestValue4", firstElement.Value.(*CacheNode).value)
	assert.Equal(t, "TestValue1", secondElement.Value.(*CacheNode).value)
	assert.Equal(t, "TestValue3", thirdElement.Value.(*CacheNode).value)
}

func TestPriorityLruCache_SetExpireOldValue(t *testing.T) {
	lruCache := New(1)
	lruCache.Set("OldKey", "OldValue", 1, time.Date(1990, 0, 0, 0, 0, 0, 0, time.UTC))
	lruCache.Set("NewKey", "NewValue", 1, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, "NewValue", lruCache.keyValuePair["NewKey"].value)
	assert.Equal(t, "NewValue", lruCache.timestampQueue[0].value)
	assert.Equal(t, 1, lruCache.timestampQueue.Len())
	assert.Equal(t, "NewValue", lruCache.priorityMap[1].Front().Value.(*CacheNode).value)
}

func TestPriorityLruCache_add(t *testing.T) {
	lruCache := New(5)
	lruCache.add("testkey", "testValue", 1, time.Now())

	assert.Equal(t, "testValue", lruCache.keyValuePair["testkey"].value)
	assert.Equal(t, "testValue", lruCache.timestampQueue[0].value)
	assert.Equal(t, "testValue", lruCache.priorityMap[1].Front().Value.(*CacheNode).value)
}

func TestPriorityLruCache_evictExpiredValue(t *testing.T) {
	lruCache := New(5)
	lruCache.add("testkey", "testValue", 1, time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC))
	lruCache.evict()

	assert.Empty(t, lruCache.timestampQueue)
	assert.Empty(t, lruCache.priorityMap)
	assert.Empty(t, lruCache.keyValuePair["testKey"])
}

func TestPriorityLruCache_doNotEvictNotExpiredValue(t *testing.T) {
	lruCache := New(5)
	lruCache.add("testkey", "testValue", 1, time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC))
	lruCache.evict()

	assert.Equal(t, "testValue", lruCache.keyValuePair["testkey"].value)
	assert.Equal(t, "testValue", lruCache.timestampQueue[0].value)
	assert.Equal(t, "testValue", lruCache.priorityMap[1].Front().Value.(*CacheNode).value)
}
