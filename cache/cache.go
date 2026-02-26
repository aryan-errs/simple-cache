package cache

import (
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
	expiryHeap  heap.Interface
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
	c.head.Next = lruNode.Next
	lruNode.Next.Prev = c.head
	lruNode.Prev = nil
	lruNode.Next = nil
}

func (c *cache) Get(Key int) int {
	if currNode, ok := c.KeyMap[Key]; ok {
		c.remove(currNode)
		if time.Now().After(currNode.ExpiresAt) {
			delete(c.KeyMap, Key)
			return -1
		}
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
		heap.Push(c.expiryHeap, types.ExpiryItem{
			ExpiresAt: existingNode.ExpiresAt,
			Node:      existingNode,
		})
		return false
	}
	newNode := &types.Node{Key: Key, Value: val, ExpiresAt: time.Now().Add(c.ttlDuration)}
	c.addToMru(newNode)
	c.KeyMap[Key] = newNode
	if c.expiryHeap.Len() > 0 {

	}
	if len(c.KeyMap) > c.capacity {
		c.evictLru()
	}
	return true
}
