package main

import (
	"fmt"
	"os"
	"time"
)

type ClientMeta struct {
	Name      string
	Topic     string
	Offset    uint32
	Prefix    string
	IPC       string
	Connected bool
	LastSeen  time.Time
}

type Hashtable struct {
	Table []*Bucket
	count int
}

func (h *Hashtable) Insert(meta *ClientMeta) {
	count := float32(h.count)
	size := float32(len(h.Table))
	if count/size > 0.75 {
		h.resize()
	}
	index := hash(meta.Name, len(h.Table))
	if h.Table[index].insert(meta) {
		h.count += 1
	}
}

func (h *Hashtable) Get(name string) *ClientMeta {
	index := hash(name, len(h.Table))
	return h.Table[index].search(name)
}

func (h *Hashtable) Delete(name string) {
	index := hash(name, len(h.Table))
	if h.Table[index].search(name) != nil {
		h.Table[index].delete(name)
		h.count -= 1
	} else {
		fmt.Fprintln(os.Stderr, "delete: couldn't find client")
	}
}

type Bucket struct {
	head *bucketNode
}

func (b *Bucket) insert(meta *ClientMeta) bool {
	if b.search(meta.Name) != nil {
		fmt.Fprintln(os.Stderr, meta.Name, ": already exists")
		return false
	}
	newNode := &bucketNode{key: meta.Name, value: meta}
	newNode.next = b.head
	b.head = newNode
	return true

}

func (b *Bucket) search(name string) *ClientMeta {
	currentNode := b.head
	for currentNode != nil {
		if currentNode.key == name {
			return currentNode.value
		}
		currentNode = currentNode.next
	}
	return nil
}

func (b *Bucket) delete(name string) {
	if b.head.key == name {
		b.head = b.head.next
		return
	}
	previousNode := b.head
	for previousNode.next != nil {
		if previousNode.next.key == name {
			previousNode.next = previousNode.next.next
			return
		}
		previousNode = previousNode.next
	}
}

type bucketNode struct {
	key   string
	value *ClientMeta
	next  *bucketNode
}

// djb2 hashing algorithm
func hash(key string, size int) int {
	var hash uint32 = 5381
	for i := 0; i < len(key); i++ {
		hash = hash*33 + uint32(key[i])
	}
	return int(hash % uint32(size))
}

func (h *Hashtable) resize() {
	newT := newTable(len(h.Table) * 2)
	for i := range h.Table {
		currentNode := h.Table[i].head
		for currentNode != nil {
			index := hash(currentNode.key, len(newT))
			newT[index].insert(currentNode.value)
			currentNode = currentNode.next
		}
	}
	h.Table = newT
}

func newTable(size int) []*Bucket {
	table := make([]*Bucket, size)
	for i := range table {
		table[i] = &Bucket{}
	}
	return table
}

func NewHashtable() *Hashtable {
	return &Hashtable{Table: newTable(64)}
}
