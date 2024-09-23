package grpc_handler

import (
	"context"

	"github.com/tricong1998/go-ecom/cmd/user/internal/services"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/pb"
)

type Server struct {
	UserService services.IUserService
	pb.UnimplementedUserGrpcServer
}

func NewServer(UserService services.IUserService) *Server {
	server := Server{
		UserService: UserService,
	}
	return &server
}

func (server *Server) ReadUser(_ context.Context, input *pb.ReadUserRequest) (*pb.User, error) {
	user, err := server.UserService.ReadUser((uint)(input.GetId()))
	if err != nil {
		return nil, err
	}
	return &pb.User{
		Username: user.Username,
		FullName: user.FullName,
		Id:       uint64(user.ID),
	}, nil
}
