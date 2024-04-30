package main

import (
	"context"
	"log"
	"net"

	"user"

	"google.golang.org/grpc"
)

type userServiceServer struct{}

func (s *userServiceServer) AddUser(ctx context.Context, user *user.User) (*user.UserResponse, error) {
	// Logic to add user (e.g., store in database)
	return &user.UserResponse{Id: user.Id}, nil
}

func (s *userServiceServer) GetUser(ctx context.Context, req *user.UserRequest) (*user.User, error) {
	// Logic to retrieve user by ID (e.g., fetch from database)
	return &user.User{Id: req.Id, Name: "John Doe", Email: "john@example.com"}, nil
}

func (s *userServiceServer) ListUsers(ctx context.Context, req *user.Empty) (*user.Users, error) {
	// Logic to list all users (e.g., fetch from database)
	users := []*user.User{
		{Id: 1, Name: "John Doe", Email: "john@example.com"},
		{Id: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	return &user.Users{Users: users}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &userServiceServer{})

	log.Println("Starting gRPC server on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
