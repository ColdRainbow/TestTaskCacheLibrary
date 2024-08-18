package cache

import (
	"lrucache/ptr"
	"strconv"
	"sync"
	"testing"
)

func TestAddFront(t *testing.T) {
	// Create a new node
	node := newNode(ptr.To("value"), "key")

	// Create a new head node
	head := newNode(ptr.To("head"), "head")

	// Add the node to the front of the list
	addToFront(head, node)

	// Check if the node was added to the front
	if head.next != node {
		t.Errorf("Expected node to be added to the front")
	}
}

func TestAddBack(t *testing.T) {
	// Create a new node
	node := newNode[string](ptr.To("value"), "key")

	// Create a new head node
	head := newNode[string](ptr.To("head"), "head")

	// Add the node to the back of the list
	addToBack(head, node)

	// Check if the node was added to the back
	if head.prev != node {
		t.Errorf("Expected node to be added to the back")
	}
}

func TestRemoveNode(t *testing.T) {
	// Create a middle new node
	nodeMiddle := newNode[string](ptr.To("value"), "key")

	// Create a new head node
	nodeBeg := newNode[string](ptr.To("value"), "key")

	// Create a new back node
	nodeBack := newNode[string](ptr.To("value"), "key")

	// Connect the nodes
	addToFront(nodeBeg, nodeMiddle)
	addToBack(nodeBeg, nodeBack)

	// Remove the node
	removeNode(nodeMiddle)

	// Check if the node was removed
	if nodeBeg.next == nodeMiddle || nodeBack.prev == nodeMiddle {
		t.Errorf("Expected node to be removed")
	}
}

// Successfully got node
func TestGetSuccess(t *testing.T) {
	// Create a new cache
	cache := NewLRUCache[string](2)

	// Add a new node
	node := newNode[string](ptr.To("value"), "key")
	cache.Put(node.key, node.value)

	// Get the node
	v, ok := cache.Get(node.key)

	// Check if the node was found
	if *v != "value" || !ok {
		t.Errorf("Expected value to be 'value'")
	}
}

// Failed to get node
func TestGetFailure(t *testing.T) {
	// Create a new cache
	cache := NewLRUCache[string](2)

	// Get the node
	v, ok := cache.Get("key2")

	// Check if the node was found
	if v != nil || ok {
		t.Errorf("Expected value to be 'value'")
	}
}

func TestPutConcurrent(t *testing.T) {
	// Create a new cache
	cache := NewLRUCache[string](2)

	// Put the node
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			cache.Put("key"+strconv.Itoa(i), ptr.To("value"+strconv.Itoa(i)))
			wg.Done()
		}()
	}
	wg.Wait()

	// Check capacity to make sure the nodes are being removed
	if len(cache.cache) != cache.capacity {
		t.Errorf("Expected cache to have less than 1000 elements")
	}
}

func TestPut(t *testing.T) {
	// Create a new cache
	cache := NewLRUCache[string](2)

	// Put the node
	cache.Put("key", ptr.To("value"))

	// Check if the node was added
	node := cache.cache["key"]
	if *node.value != "value" || node.key != "key" {
		t.Errorf("Expected node to be added")
	}

	// Check if the node was added
	if *cache.head.next.value != "value" || cache.head.next.key != "key" {
		t.Errorf("Expected node to be added to the front")
	}
}
