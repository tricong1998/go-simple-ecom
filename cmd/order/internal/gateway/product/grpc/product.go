package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/tricong1998/go-ecom/cmd/product/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IProductGateway interface {
	Get(ctx context.Context, productId uint) (*pb.ReadProductResponse, error)
	UpdateProductQuantity(ctx context.Context, productId uint, quantity uint) (bool, error)
}

type ProductGateway struct {
	host string
	port string
}

func New(host string, port string) *ProductGateway {
	return &ProductGateway{host, port}
}

func (g *ProductGateway) Get(ctx context.Context, productId uint) (*pb.ReadProductResponse, error) {
	address := fmt.Sprintf("%s:%s", g.host, g.port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewProductGrpcClient(conn)
	resp, err := client.ReadProduct(ctx, &pb.ReadProductRequest{Id: (uint64)(productId)})
	if err != nil {
		log.Println("Error getting product:", err)
		return nil, err
	}
	return resp, nil
}

func (g *ProductGateway) UpdateProductQuantity(ctx context.Context, productId uint, quantity uint) (bool, error) {
	address := fmt.Sprintf("%s:%s", g.host, g.port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := pb.NewProductGrpcClient(conn)
	resp, err := client.UpdateProductQuantity(ctx, &pb.UpdateProductQuantityRequest{ProductId: (uint64)(productId), Quantity: (uint64)(quantity)})
	if err != nil {
		log.Println("Error updating product quantity:", err)
		return false, err
	}
	return resp.Success, nil
}
