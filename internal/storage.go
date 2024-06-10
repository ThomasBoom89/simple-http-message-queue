package internal

import (
	"bufio"
	"os"
)

type Storage struct {
	queue Queue[string]
}

func NewStorage(queue Queue[string]) *Storage {
	return &Storage{
		queue: queue,
	}
}

func (S *Storage) clear() {
	err := os.WriteFile("./data", []byte(""), 0644)
	if err != nil {
		panic(err)
	}
}

func (S *Storage) Save() {
	S.clear()

	file, err := os.OpenFile("./data", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for !S.queue.IsEmpty() {
		value, _ := S.queue.Dequeue()
		//err := os.OpenFile("./data", []byte(value+"\n"), 0644)
		_, err = file.WriteString(value + "\n")
		if err != nil {
			panic(err)
		}
	}
}

func (S *Storage) Load() {
	file, err := os.OpenFile("./data", os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		S.queue.Enqueue(text)
	}
}
