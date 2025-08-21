package internal

import (
	"strconv"
	"testing"
)

func TestQueue(t *testing.T) {
	concreteQueues := []Queue[string]{NewLinkedListQueue[string](), NewSliceQueue[string]()}
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

// 1.24            32202709	        35.99 ns/op
// 1.25+greenteagc 45195168	        23.80 ns/op
// optimized       69901594	        16.55 ns/op
func BenchmarkSliceQueue(b *testing.B) {
	queue := NewSliceQueue[string]()
	for i := 0; i < b.N; i++ {
		val := strconv.Itoa(i)
		queue.Enqueue(val)
		queue.IsEmpty()
	}
	for i := 0; i < b.N; i++ {
		queue.Dequeue()
	}
}

// 1.24            21922179	        52.32 ns/op
// 1.25+greenteagc 18488520	        57.76 ns/op
func BenchmarkLinkedListQueue(b *testing.B) {
	queue := NewLinkedListQueue[string]()
	for i := 0; i < b.N; i++ {
		val := strconv.Itoa(i)
		queue.Enqueue(val)
		queue.IsEmpty()
	}
	for i := 0; i < b.N; i++ {
		queue.Dequeue()
	}
}

func BenchmarkMixedFeelingsLinkedList(b *testing.B) {
	queue := NewLinkedListQueue[string]()
	for i := 0; i < b.N; i++ {
		val := strconv.Itoa(i)
		queue.Enqueue(val)
		queue.Dequeue()
	}
}

func BenchmarkMixedFeelingsSlice(b *testing.B) {
	queue := NewSliceQueue[string]()
	for i := 0; i < b.N; i++ {
		val := strconv.Itoa(i)
		queue.Enqueue(val)
		queue.Dequeue()
	}
}
