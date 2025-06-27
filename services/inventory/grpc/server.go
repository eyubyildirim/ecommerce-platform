package grpc

import (
	"context"
	pb "ecommerce-platform/pkg/grpc/inventory"
	"ecommerce-platform/services/inventory/service"
)

type Server struct {
	pb.UnimplementedInventoryServiceServer
	service service.InventoryService
}

func NewInventoryGRPCServer(service service.InventoryService) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) GetProductInfo(ctx context.Context, req *pb.GetProductInfoRequest) (*pb.GetProductInfoResponse, error) {
	// Use the existing inventory service to get data from the database.
	products, err := s.service.GetProductsByIDs(ctx, req.ProductIds)
	if err != nil {
		return nil, err
	}

	var productInfos []*pb.ProductInfo
	for _, p := range products {
		productInfos = append(productInfos, &pb.ProductInfo{
			Id:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		})
	}

	return &pb.GetProductInfoResponse{Products: productInfos}, nil
}
