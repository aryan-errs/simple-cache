package cache

import (
	"container/heap"
	"time"
)

// TODO: add ttl to live strategy
type node struct {
	key       int
	value     int
	prev      *node
	next      *node
	expiresAt time.Time
}

type cache struct {
	capacity    int
	keyMap      map[int]*node
	head        *node
	tail        *node
	ttlDuration time.Duration
	expiryHeap  heap.Interface
}

func CreateCache(capacity int) *cache {
	c := &cache{
		capacity:    capacity,
		keyMap:      make(map[int]*node),
		head:        &node{},
		tail:        &node{},
		ttlDuration: time.Duration(10 * time.Second),
	}
	c.head.next = c.tail
	c.tail.prev = c.head
	return c
}

func (c *cache) remove(currNode *node) {
	prevNode := currNode.prev
	nextNode := currNode.next

	prevNode.next = nextNode
	nextNode.prev = prevNode

	currNode.prev = nil
	currNode.next = nil
}

func (c *cache) addToMru(newNode *node) {
	lastRealNode := c.tail.prev
	lastRealNode.next = newNode
	newNode.prev = lastRealNode
	c.tail.prev = newNode
	newNode.next = c.tail
}

func (c *cache) evictLru() {
	lruNode := c.head.next
	delete(c.keyMap, lruNode.key)
	c.head.next = lruNode.next
	lruNode.next.prev = c.head
	lruNode.prev = nil
	lruNode.next = nil
}

func (c *cache) Get(key int) int {
	if currNode, ok := c.keyMap[key]; ok {
		c.remove(currNode)
		if time.Now().After(currNode.expiresAt) {
			delete(c.keyMap, key)
			return -1
		}
		c.addToMru(currNode)
		return currNode.value
	}
	return -1
}

func (c *cache) Put(key int, val int) bool {
	if existingNode, ok := c.keyMap[key]; ok {
		c.remove(existingNode)
		existingNode.value = val
		existingNode.expiresAt = time.Now().Add(c.ttlDuration)
		c.addToMru(existingNode)
		return false
	}
	newNode := &node{key: key, value: val, expiresAt: time.Now().Add(c.ttlDuration)}
	c.addToMru(newNode)
	c.keyMap[key] = newNode
	if len(c.keyMap) > c.capacity {
		c.evictLru()
	}
	return true
}
