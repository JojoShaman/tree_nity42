package main

import (
	"fmt"
	"os"
)

const TableSize = 10

type Hashtable struct {
	Table [TableSize]*Bucket
}

func (h *Hashtable) Insert(key string) {
	index := hash(key)
	h.Table[index].insert(key)
}

func (h *Hashtable) Search(key string) bool {
	index := hash(key)
	return h.Table[index].search(key)
}

func (h *Hashtable) Delete(key string) {
	index := hash(key)
	h.Table[index].delete(key)
}

type Bucket struct {
	head *bucketNode
}

func (b *Bucket) insert(k string) {
	if !b.search(k) {

		newNode := &bucketNode{key: k}
		newNode.next = b.head
		b.head = newNode
	} else {
		fmt.Fprintln(os.Stderr, k, ": already exists")
	}
}

func (b *Bucket) search(k string) bool {
	currentNode := b.head
	for currentNode != nil {
		if currentNode.key == k {
			return true
		}
		currentNode = currentNode.next
	}
	return false
}

func (b *Bucket) delete(k string) {
	if b.head.key == k {
		b.head = b.head.next
		return
	}
	previousNode := b.head
	for previousNode.next != nil {
		if previousNode.next.key == k {
			previousNode.next = previousNode.next.next
		}
		previousNode = previousNode.next
	}
}

type bucketNode struct {
	key  string
	next *bucketNode
}

func hash(key string) int {
	sum := 0
	for _, c := range key {
		sum += int(c)
	}
	return sum % TableSize
}

func Init() *Hashtable {
	result := &Hashtable{}
	for i := range result.Table {
		result.Table[i] = &Bucket{}
	}
	return result
}
