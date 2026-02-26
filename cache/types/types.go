package types

import (
	"time"
)

type Node struct {
	Key       int
	Value     int
	Prev      *Node
	Next      *Node
	ExpiresAt time.Time
}

type ExpiryItem struct {
	ExpiresAt time.Time
	Node      *Node
}
