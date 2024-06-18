package internal

import (
	"errors"
	"github.com/gofiber/contrib/websocket"
)

type TopicManager struct {
	messageQueue    map[Topic]Queue[Message]
	connectionQueue map[Topic]Queue[*websocket.Conn]
	connections     map[Topic]map[*websocket.Conn]bool
	topics          []Topic
}

type Topic string
type Message []byte

func NewTopicManager() *TopicManager {
	return &TopicManager{
		messageQueue:    make(map[Topic]Queue[Message]),
		connectionQueue: make(map[Topic]Queue[*websocket.Conn]),
		connections:     map[Topic]map[*websocket.Conn]bool{},
	}
}

func (T *TopicManager) GetNextMessage(topic Topic) (Message, error) {
	if _, ok := T.messageQueue[topic]; !ok {
		return nil, errors.New("topic not found")
	}

	message, err := T.messageQueue[topic].Dequeue()
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (T *TopicManager) AddMessage(topic Topic, message Message) {
	if _, ok := T.messageQueue[topic]; !ok {
		T.messageQueue[topic] = NewLinkedListQueue[Message]()
		T.topics = append(T.topics, topic)
	}
	T.messageQueue[topic].Enqueue(message)
}

func (T *TopicManager) AddConnection(topic Topic, conn *websocket.Conn) {
	if _, ok := T.connections[topic]; !ok {
		T.connections[topic] = map[*websocket.Conn]bool{conn: true}
	} else {
		T.connections[topic][conn] = true
	}
}

func (T *TopicManager) RemoveConnection(topic Topic, conn *websocket.Conn) {
	if _, ok := T.connections[topic]; !ok {
		return
	}
	delete(T.connections[topic], conn)
}

func (T *TopicManager) GetTopics() []Topic {
	return T.topics
}

func (T *TopicManager) GetNextConnection(topic Topic) (*websocket.Conn, error) {
	if T.topicNotReady(topic) {
		return nil, errors.New("topic not found or no data")
	}

	connection, err := T.connectionQueue[topic].Dequeue()
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func (T *TopicManager) ConnectionExists(topic Topic, connection *websocket.Conn) bool {
	if _, ok := T.connections[topic]; !ok {
		return false
	}
	if _, ok := T.connections[topic][connection]; !ok {
		return false
	}

	return true
}

func (T *TopicManager) AddConnectionToQueue(topic Topic, connection *websocket.Conn) {
	if _, ok := T.connectionQueue[topic]; !ok {
		T.connectionQueue[topic] = NewLinkedListQueue[*websocket.Conn]()
	}
	T.connectionQueue[topic].Enqueue(connection)
}

func (T *TopicManager) topicNotReady(topic Topic) bool {
	if _, ok := T.connections[topic]; !ok {
		return true
	}
	if _, ok := T.connectionQueue[topic]; !ok {
		return true
	}
	if _, ok := T.messageQueue[topic]; !ok {
		return true
	}
	return T.messageQueue[topic].IsEmpty() || T.connectionQueue[topic].IsEmpty() || len(T.connections[topic]) == 0
}
