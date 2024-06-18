package internal

import (
	"bufio"
	"os"
)

const DataDir = "./data/"
const FileExcluded = ".gitkeep"

type Storage struct {
	topicManager *TopicManager
}

func NewStorage(topicManager *TopicManager) *Storage {
	return &Storage{
		topicManager: topicManager,
	}
}

func (S *Storage) createStoreDirectory() {
	_, err := os.ReadDir(DataDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(DataDir, 0755)
			if err != nil {
				panic(err)
			}
		} else if !os.IsExist(err) {
			panic(err)
		}
	}
}

func (S *Storage) clear() {
	dirEntries, err := os.ReadDir(DataDir)
	if err != nil {
		panic(err)
	}
	for _, entry := range dirEntries {
		if entry.Name() == FileExcluded {
			continue
		}
		os.Remove(DataDir + entry.Name())
	}
}

func (S *Storage) Save() {
	S.clear()
	S.createStoreDirectory()

	for _, topic := range S.topicManager.GetTopics() {
		S.saveTopic(topic)
	}
}

func (S *Storage) saveTopic(topic Topic) {
	file, err := os.OpenFile(string(DataDir+topic), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for {
		message, err := S.topicManager.GetNextMessage(topic)
		if err != nil {
			return
		}
		_, err = file.Write(message)
		if err != nil {
			panic(err)
		}
		// Line feed
		_, _ = file.Write([]byte{10})
	}
}

func (S *Storage) Load() {
	dirEntries, err := os.ReadDir(DataDir)
	if err != nil {
		panic(err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.Name() == FileExcluded {
			continue
		}
		S.loadTopic(dirEntry)
	}
}

func (S *Storage) loadTopic(dirEntry os.DirEntry) {
	topic := dirEntry.Name()
	file, err := os.OpenFile(DataDir+topic, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Bytes()
		if len(text) <= 0 {
			continue
		}
		S.topicManager.AddMessage(Topic(topic), text)
	}
}
