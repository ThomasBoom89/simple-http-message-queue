package grpc

import (
	"context"
	"github.com/ThomasBoom89/simple-http-message-queue/internal"
	"github.com/gofiber/fiber/v2/log"
)

type Server struct {
	topicManager *internal.TopicManager
}

func NewServer(topicManager *internal.TopicManager) *Server {
	return &Server{topicManager: topicManager}
}

func (S *Server) Publish(ctx context.Context, request *PublishRequest) (*Empty, error) {
	log.Debug(request.GetMessage())
	S.topicManager.AddMessage(internal.Topic(request.GetTopic()), request.GetMessage())

	return &Empty{}, nil
}

func (S *Server) Subscribe(ctx context.Context, request *SubscribeRequest) (*Response, error) {
	message, err := S.topicManager.GetNextMessage(internal.Topic(request.GetTopic()))
	if err != nil {
		return &Response{Message: []byte{}}, nil
	}

	return &Response{Message: message}, nil
}

func (S *Server) mustEmbedUnimplementedMessageBrokerServer() {
}
