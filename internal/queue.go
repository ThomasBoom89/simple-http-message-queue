package internal

import (
	"errors"
	"github.com/gofiber/contrib/websocket"
)

type Queue[T string | *websocket.Conn] interface {
	IsEmpty() bool
	Enqueue(T)
	Dequeue() (T, error)
}

type SliceQueue[T string | *websocket.Conn] struct {
	data []T
}

func NewSliceQueue[T string | *websocket.Conn]() *SliceQueue[T] {
	return &SliceQueue[T]{data: nil}
}

func (S *SliceQueue[T]) IsEmpty() bool {
	return len(S.data) == 0
}

func (S *SliceQueue[T]) Enqueue(element T) {
	S.data = append(S.data, element)
}

func (S *SliceQueue[T]) Dequeue() (T, error) {
	if S.IsEmpty() {
		return *new(T), errors.New("messageQueue is empty")
	}
	element := S.data[0]
	S.data = S.data[1:]

	return element, nil
}

type LinkedList[T string | *websocket.Conn] struct {
	head *Node[T]
	tail *Node[T]
}

type Node[T string | *websocket.Conn] struct {
	data T
	next *Node[T]
}

type LinkedListQueue[T string | *websocket.Conn] struct {
	data *LinkedList[T]
}

func NewLinkedListQueue[T string | *websocket.Conn]() *LinkedListQueue[T] {
	return &LinkedListQueue[T]{data: &LinkedList[T]{head: nil, tail: nil}}
}

func (L *LinkedListQueue[T]) IsEmpty() bool {
	return L.data.head == nil
}

func (L *LinkedListQueue[T]) Enqueue(element T) {
	node := &Node[T]{data: element}
	if L.IsEmpty() {
		L.data.head = node
	}
	tail := L.data.tail
	if tail != nil {
		tail.next = node
	}
	L.data.tail = node
}

func (L *LinkedListQueue[T]) Dequeue() (T, error) {
	if L.IsEmpty() {
		return *new(T), errors.New("messageQueue is empty")
	}
	node := L.data.head
	if node.next == nil {
		L.data.tail = nil
	}
	L.data.head = node.next
	node.next = nil

	return node.data, nil
}
