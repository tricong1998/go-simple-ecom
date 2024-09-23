package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/tricong1998/go-ecom/cmd/user/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	// client pb.UserGrpcClient
	host string
	port string
}

func New(host string, port string) *Gateway {
	return &Gateway{host, port}
}

func (g *Gateway) Get(ctx context.Context, userId uint) (*pb.User, error) {
	address := fmt.Sprintf("%s:%s", g.host, g.port)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	fmt.Println("error")

	client := pb.NewUserGrpcClient(conn)
	fmt.Println("error2")
	resp, err := client.ReadUser(ctx, &pb.ReadUserRequest{Id: (uint64)(userId)})
	if err != nil {
		log.Println("Error getting user:", err)
		return nil, err
	}
	return resp, nil
}
