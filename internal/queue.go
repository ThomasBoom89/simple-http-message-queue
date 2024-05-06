package internal

import "errors"

type Queue interface {
	IsEmpty() bool
	Enqueue(string)
	Dequeue() (string, error)
}

type SliceQueue struct {
	data []string
}

func NewSliceQueue() *SliceQueue {
	return &SliceQueue{data: nil}
}

func (S *SliceQueue) IsEmpty() bool {
	return len(S.data) == 0
}

func (S *SliceQueue) Enqueue(element string) {
	S.data = append(S.data, element)
}

func (S *SliceQueue) Dequeue() (string, error) {
	if S.IsEmpty() {
		return "", errors.New("queue is empty")
	}
	element := S.data[0]
	S.data = S.data[1:]

	return element, nil
}

type LinkedList struct {
	head *Node
	tail *Node
}

type Node struct {
	data string
	next *Node
}

type LinkedListQueue struct {
	data *LinkedList
}

func NewLinkedListQueue() *LinkedListQueue {
	return &LinkedListQueue{data: &LinkedList{head: nil, tail: nil}}
}

func (L *LinkedListQueue) IsEmpty() bool {
	return L.data.head == nil
}

func (L *LinkedListQueue) Enqueue(element string) {
	node := &Node{data: element}
	if L.IsEmpty() {
		L.data.head = node
	}
	tail := L.data.tail
	if tail != nil {
		tail.next = node
	}
	L.data.tail = node
}

func (L *LinkedListQueue) Dequeue() (string, error) {
	if L.IsEmpty() {
		return "", errors.New("queue is empty")
	}
	node := L.data.head
	if node.next == nil {
		L.data.tail = nil
	}
	L.data.head = node.next
	node.next = nil

	return node.data, nil
}
