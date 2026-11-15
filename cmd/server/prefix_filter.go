package main

type PrefixFilter struct {
	root *Node
}

type Node struct {
	children  map[byte]*Node
	consumers []*Consumer
}

func newNode() *Node {
	return &Node{
		children: make(map[byte]*Node),
	}
}

func (f *PrefixFilter) Insert(prefix string, consumer *Consumer) {
	curr := f.root
	for i := 0; i < len(prefix); i++ {
		child, exists := curr.children[prefix[i]]
		if !exists {
			child = newNode()
			curr.children[prefix[i]] = child
		}
		curr = child
	}
	curr.consumers = append(curr.consumers, consumer)
}

func (f *PrefixFilter) Match(key string) []*Consumer {
	var match []*Consumer
	curr := f.root
	match = append(match, curr.consumers...)
	for i := 0; i < len(key); i++ {
		child, exists := curr.children[key[i]]
		if !exists {
			break
		}
		curr = child
		match = append(match, curr.consumers...)
	}
	return match
}

func (f *PrefixFilter) Remove(prefix string, consumer *Consumer) {
	curr := f.root
	for i := 0; i < len(prefix); i++ {
		child, exists := curr.children[prefix[i]]
		if !exists {
			return
		}
		curr = child
	}
	for i := 0; i < len(curr.consumers); i++ {
		if curr.consumers[i] == consumer {
			curr.consumers = append(
				curr.consumers[:i], curr.consumers[i+1:]...)
			return
		}
	}
}
