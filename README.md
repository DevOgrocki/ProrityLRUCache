# ProrityLRUCache
A cache with expiry based on priority, timestamp, and recency

Calling `PriorityLruCache.New(N)` creates a new cache of `N` size.

Calling `PriorityLruCache.Set()`, allows you to set a key, value, priority, and expiry.

Calling `PriorityLruCache.Get()` will return that value, Get will never return an expired value.

If you attempt to insert more values into the cache then `N` The cache will replace values in this order:
1. Remove expired values
2. Remove the lowest Priority values and in case of a priority tie remove the least recently used.

Example code:

```
	// Create cache of with a max size of 3
	lruCache := New(3)
	// [ _, _, _ ]
	
	// Add a key value pair of A:B
	lruCache.Set("A", "B", 1, time.Date(2223, 0, 0, 0, 0, 0, 0, time.UTC))
	// [ A:B, _, _ ]

	// Add a key value pair of C:D
	lruCache.Set("C", "D", 1, time.Date(2223, 0, 0, 0, 0, 0, 0, time.UTC))
	// [ A:B, C:D, _ ]

	// Add a key value pair of E:F
	lruCache.Set("E", "F", 2, time.Date(2223, 0, 0, 0, 0, 0, 0, time.UTC))
	// [ A:B P1, C:D P1, E:F P2 ]

	// Add a key value pair of G:H, all values are not expired so replace the lowest priority value E:F
	lruCache.Set("G", "H", 1, time.Date(2223, 0, 0, 0, 0, 0, 0, time.UTC))
	// [ A:B P1, C:D P1, G:H P1 ]

	// Add a key value pair of I:J, All values are not expired, all have the same priority, remove the least recently used one A:B
	lruCache.Set("I", "J", 1, time.Date(1990, 0, 0, 0, 0, 0, 0, time.UTC))
	// [ I:J P1, C:D P1, G:H P1 ]

	lruCache.Get("I") //return ("", false) because I:J is expired
	
	// Add a key value pair of K:L, the I:J entry is expired, replace it.
	lruCache.Set("K", "L", 1, time.Date(2223, 0, 0, 0, 0, 0, 0, time.UTC))
	// [ K:L P1, C:D P1, G:H P1 ]
	
	lruCache.Get("C") //return ("D", true) and makes it the most recently used

	// Add a key value pair of M:N, Replace the least recently used key G:H
	lruCache.Set("K", "L", 1, time.Date(2223, 0, 0, 0, 0, 0, 0, time.UTC))
	// [ K:L P1, C:D P1, K:L P1 ]
```

## Contribution

Please 