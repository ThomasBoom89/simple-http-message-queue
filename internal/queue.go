package internal

import "errors"

type Queue struct {
	data []string
}

func NewQueue() *Queue {
	return &Queue{data: nil}
}

func (Q *Queue) IsEmpty() bool {
	return len(Q.data) == 0
}

func (Q *Queue) Enqueue(element string) {
	Q.data = append(Q.data, element)
}

func (Q *Queue) Dequeue() (string, error) {
	if Q.IsEmpty() {
		return "", errors.New("queue is empty")
	}
	element := Q.data[0]
	Q.data = Q.data[1:]

	return element, nil
}
