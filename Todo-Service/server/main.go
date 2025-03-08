package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	pb "github.com/valdimir-makarov/Go-backend-Engineering/Todo-Service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedTodoServiceServer

	todos  map[int64]*pb.TodoItem
	mu     sync.Mutex
	nextID int64 // Counter for generating unique numeric IDs

}

func (s *server) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.TodoItem, error) {
	if req == nil || req.Title == "" || req.Description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Title and Description are required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	id := atomic.AddInt64(&s.nextID, 1)
	todo := &pb.TodoItem{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false, // Default completion status is false
	}
	s.todos[id] = todo
	return todo, nil
}

func (s *server) ListTodosResponse(ctx context.Context, req *pb.CreateTodoRequest) (*pb.ListTodosResponse, error) {

	todos := make([]*pb.TodoItem, 0, len(s.todos))

	for key, value := range s.todos {
		if value != nil { // Ensure the pointer is not nil

			todos = append(todos, value)

			fmt.Printf("Key: %d, Title: %s, Description: %s, Completed: %t\n",
				key, value.Title, value.Description, value.Completed)

		} else {
			fmt.Printf("Key: %d, Value: <nil>\n", key)
		}
	}
	return &pb.ListTodosResponse{
		Todos: todos,
	}, nil
}
