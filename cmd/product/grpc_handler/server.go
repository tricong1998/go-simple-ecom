package grpc_handler

import (
	"context"

	"github.com/tricong1998/go-ecom/cmd/product/internal/services"
	"github.com/tricong1998/go-ecom/cmd/product/pkg/pb"
)

type Server struct {
	ProductService services.IProductService
	pb.UnimplementedProductGrpcServer
}

func NewServer(ProductService services.IProductService) *Server {
	server := Server{
		ProductService: ProductService,
	}
	return &server
}

func (server *Server) ReadProduct(_ context.Context, input *pb.ReadProductRequest) (*pb.ReadProductResponse, error) {
	product, err := server.ProductService.ReadProduct((uint)(input.GetId()))
	if err != nil {
		return nil, err
	}
	return &pb.ReadProductResponse{
		Product: &pb.Product{
			Name:     product.Name,
			Price:    uint64(product.Price),
			Quantity: uint64(product.Quantity),
			Id:       uint64(product.ID),
		},
	}, nil
}

func (server *Server) UpdateProductQuantity(_ context.Context, input *pb.UpdateProductQuantityRequest) (*pb.UpdateProductQuantityResponse, error) {
	success, err := server.ProductService.UpdateProductQuantity((uint)(input.GetProductId()), (uint)(input.GetQuantity()))
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProductQuantityResponse{
		Success: success,
	}, nil
}
