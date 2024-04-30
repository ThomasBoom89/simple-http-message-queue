package main

import (
	"fmt"
	"github.com/ThomasBoom89/simple-http-message-queue/internal"
	"time"
)

func main() {
	queue := internal.NewQueue()

	queue.Enqueue("a")
	queue.Enqueue("b")
	queue.Enqueue("c")

	for {
		//if !queue.IsEmpty() {
		fmt.Println(queue.Dequeue())
		//}
		time.Sleep(2 * time.Second)
	}
}
