	package cache

	import (
		expiryheap "aryan-errs/simple-cache/cache/expiryHeap"
		"aryan-errs/simple-cache/cache/types"
		"container/heap"
		"time"
	)

	type cache struct {
		capacity    int
		KeyMap      map[int]*types.Node
		head        *types.Node
		tail        *types.Node
		ttlDuration time.Duration
		expiryHeap  *expiryheap.ExpiryHeap
	}

	func CreateCache(capacity int) *cache {
		c := &cache{
			capacity:    capacity,
			KeyMap:      make(map[int]*types.Node),
			head:        &types.Node{},
			tail:        &types.Node{},
			ttlDuration: time.Duration(10 * time.Second),
		}
		c.head.Next = c.tail
		c.tail.Prev = c.head
		c.expiryHeap = &expiryheap.ExpiryHeap{}
		heap.Init(c.expiryHeap)
		return c
	}

	func (c *cache) remove(currNode *types.Node) {
		PrevNode := currNode.Prev
		NextNode := currNode.Next

		PrevNode.Next = NextNode
		NextNode.Prev = PrevNode

		currNode.Prev = nil
		currNode.Next = nil
	}

	func (c *cache) addToMru(newNode *types.Node) {
		lastRealNode := c.tail.Prev
		lastRealNode.Next = newNode
		newNode.Prev = lastRealNode
		c.tail.Prev = newNode
		newNode.Next = c.tail
	}

	func (c *cache) evictLru() {
		lruNode := c.head.Next
		delete(c.KeyMap, lruNode.Key)
		c.remove(lruNode)
	}

	func (c *cache) evict() {
		for len(c.KeyMap) > c.capacity && c.expiryHeap.Len() > 0 {
			expiredItem := c.expiryHeap.Peak()
			if expiredItem.ExpiresAt.After(time.Now()) {
				break
			}
			heap.Pop(c.expiryHeap)
			node := expiredItem.Node
			if node.ExpiresAt.Equal(expiredItem.ExpiresAt) {
				c.remove(node)
				delete(c.KeyMap, node.Key)
			}
		}
		if len(c.KeyMap) > c.capacity {
			c.evictLru()
		}
	}

	func (c *cache) Get(Key int) int {
		if currNode, ok := c.KeyMap[Key]; ok {
			if time.Now().After(currNode.ExpiresAt) {
				c.remove(currNode)
				delete(c.KeyMap, Key)
				return -1
			}
			c.remove(currNode)
			c.addToMru(currNode)
			return currNode.Value
		}
		return -1
	}

	func (c *cache) Put(Key int, val int) bool {
		if existingNode, ok := c.KeyMap[Key]; ok {
			c.remove(existingNode)
			existingNode.Value = val
			existingNode.ExpiresAt = time.Now().Add(c.ttlDuration)
			c.addToMru(existingNode)
			heap.Push(c.expiryHeap, &types.ExpiryItem{
				ExpiresAt: existingNode.ExpiresAt,
				Node:      existingNode,
			})
			return false
		}
		newNode := &types.Node{Key: Key, Value: val, ExpiresAt: time.Now().Add(c.ttlDuration)}
		c.addToMru(newNode)
		c.KeyMap[Key] = newNode
		heap.Push(c.expiryHeap, &types.ExpiryItem{
			ExpiresAt: newNode.ExpiresAt,
			Node:      newNode,
		})
		if len(c.KeyMap) > c.capacity {
			c.evict()
		}
		return true
	}
