package internal

import (
	"errors"
	"github.com/gofiber/contrib/websocket"
)

type TopicManager struct {
	data       map[Topic]Queue[Message]
	connection map[Topic]Queue[*websocket.Conn]
}

type Topic string
type Message []byte

func NewTopicManager() *TopicManager {
	return &TopicManager{
		data:       make(map[Topic]Queue[Message]),
		connection: make(map[Topic]Queue[*websocket.Conn]),
	}
}

func (T *TopicManager) GetNextMessage(topic Topic) (Message, error) {
	if _, ok := T.data[topic]; !ok {
		return nil, errors.New("topic not found")
	}

	message, err := T.data[topic].Dequeue()
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (T *TopicManager) AddMessage(topic Topic, message Message) {
	if _, ok := T.data[topic]; !ok {
		T.data[topic] = NewLinkedListQueue[Message]()
	}
	T.data[topic].Enqueue(message)
}

func (T *TopicManager) AddConnection(topic Topic, conn *websocket.Conn) {

}

func (T *TopicManager) RemoveConnection(topic Topic, conn *websocket.Conn) {
	// todo: this will not work, we do not know topic outside, maybe we have to scan for connection inside all topic
}

//func (T *TopicManager) GetTopics() []Topic {}
