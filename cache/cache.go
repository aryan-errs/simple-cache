package cache

import (
	"time"
)

// Requirements of a cache (Since I am implementing LRU cache, i have requirements based on that)
// 1. user should be able to store key, value, so that would be a put operation
// 2. user should be able to get the value using the key
// 3. when a user gets the a key, that key would now become the most frequently used one
// 4.I would also need to implement eviction mechanism that remove the LRU key
// 5. On last thing that can be added is tha time to live (TTL) for each key. (Need to look into the strategies for this)

type node struct {
	key   int
	value int
	prev  *node
	next  *node
}

type cache struct {
	capacity int
	keyMap   map[int]*node
	head     *node
	tail     *node
	ttl      time.Time
}

func CreateCache(capacity int) *cache {
	c := &cache{
		capacity: capacity,
		keyMap:   make(map[int]*node),
		head:     &node{},
		tail:     &node{},
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
		c.addToMru(currNode)
		return currNode.value
	}
	return -1
}

func (c *cache) Put(key int, val int) bool {
	if existingNode, ok := c.keyMap[key]; ok {
		c.remove(existingNode)
		existingNode.value = val
		c.addToMru(existingNode)
		return false
	}
	newNode := &node{key: key, value: val}
	c.addToMru(newNode)
	c.keyMap[key] = newNode
	if len(c.keyMap) > c.capacity {
		c.evictLru()
	}
	return true
}
