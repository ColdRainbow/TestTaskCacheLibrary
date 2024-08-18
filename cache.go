// One pacakge to incupsulate the module
package cache

import "sync"

// Uses generics to implement a node
type Node[T any] struct {
	// Uses a pointer to save memory
	value *T
	key   string
	prev  *Node[T]
	next  *Node[T]
}

func newNode[T any](value *T, key string) *Node[T] {
	// Node points to itself, our starting point
	node := &Node[T]{
		value: value,
		key:   key,
	}
	node.next = node
	node.prev = node
	return node
}

// Adds a node to the front of the list
// Not possible to export so that the user can't add nodes
func addToFront[T any](head *Node[T], node *Node[T]) {
	node.next = head.next
	node.prev = head
	head.next.prev = node
	head.next = node
}

// Adds a node to the back of the list
// Not possible to export so that the user can't add nodes
func addToBack[T any](head *Node[T], node *Node[T]) {
	node.prev = head.prev
	node.next = head
	head.prev.next = node
	head.prev = node
}

// Removes a node from the list
// Not possible to export so that the user can't remove nodes
func removeNode[T any](node *Node[T]) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// Represents the cache itself
type LRUCache[T any] struct {
	capacity int
	cache    map[string]*Node[T]
	head     *Node[T]
	mu       sync.RWMutex
}

// Initializes a new LRU cache with a given capacity
func NewLRUCache[T any](capacity int) *LRUCache[T] {
	return &LRUCache[T]{
		capacity: capacity,
		cache:    make(map[string]*Node[T]),
		head:     newNode[T](nil, ""),
	}
}

// Gets a value from the cache
// Returns only value instead of a pointer to the whole node,
// as the user doesn't need to modify the node
func (lru *LRUCache[T]) Get(key string) (*T, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if node, ok := lru.cache[key]; ok {
		removeNode(node)
		addToFront(lru.head, node)
		return node.value, true
	}
	return nil, false
}

// Puts a value into the cache
func (lru *LRUCache[T]) Put(key string, value *T) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if node, ok := lru.cache[key]; ok {
		node.value = value
		removeNode(node)
		addToFront(lru.head, node)
	}

	if len(lru.cache) >= lru.capacity {
		last := lru.head.prev
		lru.cache[last.key] = nil
		delete(lru.cache, last.key)
		removeNode(lru.head.prev)
	}

	node := newNode(value, key)
	lru.cache[key] = node
	addToFront(lru.head, node)
}
