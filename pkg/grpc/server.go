package grpc

import (
	"context"
	"fmt"
	"lucas-stellet/api-grpc-mongodb/internal/auth"
	"lucas-stellet/api-grpc-mongodb/pkg/logger"
	"lucas-stellet/api-grpc-mongodb/pkg/pb"
	"lucas-stellet/api-grpc-mongodb/pkg/services"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) SetGlobalSettings(ctx context.Context, in *pb.SetGlobalSettingsRequest) (*pb.SetGlobalSettingsResponse, error) {
	if err := auth.ValidateTokenFromMeta(ctx); err != nil {
		return nil, err
	}

	query := in.GetQuery()

	logger.Write(logger.GRPC, fmt.Sprintf("in: %s", query.String()), logger.STDOUT)

	project := services.ProjectMiddleware{}

	result, err := project.GetData(query.GetId())

	if err != nil {
		return nil, err
	}

	logger.Write(logger.GRPC, fmt.Sprintf("out: %s", result), logger.STDOUT)

	return &pb.SetGlobalSettingsResponse{
		Result: result,
	}, nil
}

func (s *server) SetProject(ctx context.Context, in *pb.SetProjectRequest) (*pb.SetProjectResponse, error) {
	if err := auth.ValidateTokenFromMeta(ctx); err != nil {
		return nil, err
	}

	query := in.GetQuery()

	logger.Write(logger.GRPC, fmt.Sprintf("in: %s", query.String()), logger.FILE)

	project := services.ProjectIdentity{}

	result, err := project.GetData(query.GetUrl(), query.GetApiVersion())

	if err != nil {
		return nil, err
	}

	logger.Write(logger.GRPC, fmt.Sprintf("out: %s", result), logger.STDOUT)

	return &pb.SetProjectResponse{
		Result: result,
	}, nil
}

func Start() {
	logger.Write(logger.INFO, "grpc server started", logger.STDOUT)

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		logger.Write(logger.FATAL, fmt.Sprintf("Error creating a TCP network on port 50051 :: %v", err), logger.FILE)
	}

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	pb.RegisterInternalGatewayServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		logger.Write(logger.FATAL, fmt.Sprintf("Error when serving gRPC server :: %v", err), logger.FILE)
	}
}
