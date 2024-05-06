package internal

import (
	"strconv"
	"testing"
)

func TestQueue(t *testing.T) {
	concreteQueues := []Queue{NewLinkedListQueue(), NewSliceQueue()}

	for _, queue := range concreteQueues {
		if queue.IsEmpty() == false {
			t.Fatalf("Queue should be empty")
		}
		if _, err := queue.Dequeue(); err == nil {
			t.Fatalf("Queue should be empty")
		}
		queue.Enqueue("a")
		queue.Enqueue("b")
		queue.Enqueue("c")
		if _, err := queue.Dequeue(); err != nil {
			t.Fatalf("Queue should not be empty")
		}
		if value, _ := queue.Dequeue(); value != "b" {
			t.Fatalf("Queue should return b")
		}

		if _, err := queue.Dequeue(); err != nil {
			t.Fatalf("Queue should not be empty")
		}

		if _, err := queue.Dequeue(); err == nil {
			t.Fatalf("Queue should be empty")
		}

		if queue.IsEmpty() == false {
			t.Fatalf("Queue should be empty")
		}
	}
}

func BenchmarkSliceQueue(b *testing.B) {
	queue := NewSliceQueue()
	for i := 0; i < b.N; i++ {
		val := strconv.Itoa(i)
		queue.Enqueue(val)
		queue.IsEmpty()
	}
	for i := 0; i < b.N; i++ {
		queue.Dequeue()
	}
}

func BenchmarkLinkedListQueue(b *testing.B) {
	queue := NewLinkedListQueue()
	for i := 0; i < b.N; i++ {
		val := strconv.Itoa(i)
		queue.Enqueue(val)
		queue.IsEmpty()
	}
	for i := 0; i < b.N; i++ {
		queue.Dequeue()
	}
}
