package pkg

import "sync"

type Node struct {
	Key       string
	Val       string
	Pre, Next *Node
}

func NewNode(key, val string) *Node {
	return &Node{
		Key: key,
		Val: val,
	}
}

type LinkedList struct {
	Size       int
	Head, Tail *Node
}

func InitLinkedList() *LinkedList {
	list := &LinkedList{
		Head: NewNode("head", ""),
		Tail: NewNode("tail", ""),
		Size: 0,
	}
	list.Head.Next = list.Tail
	list.Tail.Pre = list.Head
	return list
}
func (l *LinkedList) PushBack(node *Node) {
	node.Pre = l.Tail.Pre
	node.Next = l.Tail
	l.Tail.Pre.Next = node
	l.Tail.Pre = node
	l.Size++
}
func (l *LinkedList) PushFront(node *Node) {
	node.Pre = l.Head
	node.Next = l.Head.Next
	l.Head.Next.Pre = node
	l.Head.Next = node
	l.Size++
}
func (l *LinkedList) Remove(node *Node) {
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
	node.Pre = nil
	node.Next = nil
	l.Size--
}
func (l *LinkedList) PopFront() (ans *Node) {
	if l.Empty() {
		return
	}
	ans = l.Head.Next
	l.Remove(ans)
	return
}
func (l *LinkedList) PopBack() (ans *Node) {
	if l.Empty() {
		return
	}
	ans = l.Tail.Pre
	l.Remove(ans)
	return
}
func (l *LinkedList) Front() (ans *Node) {
	if l.Empty() {
		return
	}
	ans = l.Head.Next
	return
}
func (l LinkedList) Back() (ans *Node) {
	if l.Empty() {
		return
	}
	ans = l.Tail.Pre
	return
}
func (l *LinkedList) Empty() bool {
	return l.Size == 0 || (l.Head.Next == l.Tail && l.Tail.Pre == l.Head)
}

type LRUCache struct {
	mu       sync.RWMutex
	Capacity int
	HashMap  map[string]*Node
	Cache    *LinkedList
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		mu:       sync.RWMutex{},
		Capacity: capacity,
		HashMap:  make(map[string]*Node),
		Cache:    InitLinkedList(),
	}
}

func (l *LRUCache) Get(key string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if node, has := l.HashMap[key]; !has {
		return ""
	} else {
		l.Cache.Remove(node)
		l.Cache.PushBack(node)
		return node.Val
	}
}

func (l *LRUCache) Put(key, val string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if node, has := l.HashMap[key]; !has {
		if l.Cache.Size == l.Capacity {
			nd1 := l.Cache.PopFront()
			delete(l.HashMap, nd1.Key)
		}
		l.put(key, val)
	} else {
		l.Cache.Remove(node)
		delete(l.HashMap, node.Key)
		l.put(key, val)
	}
}
func (l *LRUCache) put(key, val string) {
	newNode := NewNode(key, val)
	l.Cache.PushBack(newNode)
	l.HashMap[key] = newNode
}
