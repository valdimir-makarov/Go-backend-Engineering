package handler

import (
	"context"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated/github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatServer struct {
	generated.UnsafeChatServiceServer
	svc *service.Service
}

func NewServer(service *service.Service) *GRPCChatServer {
	return &GRPCChatServer{
		svc: service,
	}

}
func (s *GRPCChatServer) SendMessage(ctx context.Context, req *generated.ChatMessage) (*generated.Empty, error) {

	err := s.svc.SendMessages(int(req.SenderId), int(req.ReceiverId), req.Content)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to send message yo service: %v", err)
	}
	return &generated.Empty{}, nil
}
