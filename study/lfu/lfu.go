package main

import "fmt"

type Node struct {
	Key, Val, Freq int
	Pre, Next      *Node
}

func NewNode(key, val int) *Node {
	return &Node{
		Key:  key,
		Val:  val,
		Freq: 1,
	}
}

type LinkedList struct {
	Head, Tail *Node
}

func InitLinkedList() *LinkedList {
	list := &LinkedList{
		Head: NewNode(-1, -1),
		Tail: NewNode(-1, -1),
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
}
func (l *LinkedList) PushFront(node *Node) {
	node.Pre = l.Head
	node.Next = l.Head.Next
	l.Head.Next.Pre = node
	l.Head.Next = node
}
func (l *LinkedList) Remove(node *Node) {
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
	node.Pre = nil
	node.Next = nil
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
	return l.Head.Next == l.Tail && l.Tail.Pre == l.Head
}

type LFUCache struct {
	Cache    map[int]*Node
	Freq     map[int]*LinkedList
	Capacity int
	Size     int
	MinFreq  int
}

func Constructor(cap int) LFUCache {
	return LFUCache{
		Capacity: cap,
		Cache:    make(map[int]*Node),
		Freq:     make(map[int]*LinkedList),
	}
}
func (l *LFUCache) Get(key int) int {
	if node, has := l.Cache[key]; has {
		l.incFreq(node)
		return node.Val
	}
	return -1
}

func (l *LFUCache) Put(key, val int) {
	if l.Capacity == 0 {
		return
	}
	if node, has := l.Cache[key]; has {
		node.Val = val
		l.incFreq(node)
	} else {
		if l.Size == l.Capacity {
			rm := l.Freq[l.MinFreq].PopBack()
			delete(l.Cache, rm.Key)
			l.Size--
		}
		l.put(key, val)
	}
}
func (l *LFUCache) put(key, val int) {
	newNode := NewNode(key, val)
	l.Cache[key] = newNode
	if l.Freq[1] == nil {
		l.Freq[1] = InitLinkedList()
	}
	l.Freq[1].PushFront(newNode)
	l.MinFreq = 1
	l.Size++
}
func (l *LFUCache) incFreq(node *Node) {
	freq := node.Freq
	l.Freq[node.Freq].Remove(node)
	if l.MinFreq == freq && l.Freq[l.MinFreq].Empty() {
		l.MinFreq++
		delete(l.Freq, freq)
	}
	node.Freq++
	if l.Freq[node.Freq] == nil {
		l.Freq[node.Freq] = InitLinkedList()
	}
	l.Freq[node.Freq].PushFront(node)
}
func main() {
	l := Constructor(2)
	l.Put(1, 1)
	l.Put(2, 2)
	val := l.Get(1)
	fmt.Println(val)
	l.Put(3, 3)
	val = l.Get(2)
	fmt.Println(val)
	val = l.Get(3)
	fmt.Println(val)
	l.Put(4, 4)
	val = l.Get(1)
	fmt.Println(val)
	val = l.Get(3)
	fmt.Println(val)
	val = l.Get(4)
	fmt.Println(val)
}
