package grpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rkt "github.com/Dimashey/protos-gprc/rocket/v1"
	"github.com/google/uuid"

	"github.com/Dimashey/gprc/internal/rocket"
)

type RockerService interface {
	InsertRocket(ctx context.Context, rkt rocket.Rocket) (rocket.Rocket, error)
	GetRocketByID(ctx context.Context, id string) (rocket.Rocket, error)
	DeleteRocket(ctx context.Context, id string) error
}

type Handler struct {
	RockerService RockerService
	rkt.UnimplementedRocketServiceServer
}

func New(rktService RockerService) Handler {
	return Handler{
		RockerService: rktService,
	}
}

func (h Handler) Serve() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Println("could not listen on port 50051")
		return err
	}

	grpcServer := grpc.NewServer()

	rkt.RegisterRocketServiceServer(
		grpcServer,
		&h,
	)

	if err := grpcServer.Serve(lis); err != nil {
		log.Println("failed to serve: %s\n", err)
		return err
	}

	return nil
}

func (h Handler) GetRocket(ctx context.Context, req *rkt.GetRocketRequest) (*rkt.GetRocketResponse, error) {
	log.Println("Get Rocket gRPC Endpoint Hit")

	rocket, err := h.RockerService.GetRocketByID(ctx, req.Id)
	if err != nil {
		log.Print("Failed to retrieve rocket by ID")
		return &rkt.GetRocketResponse{}, err
	}

	return &rkt.GetRocketResponse{
		Rocket: &rkt.Rocket{
			Id:   rocket.ID,
			Name: rocket.Name,
			Type: rocket.Type,
		},
	}, nil
}

func (h Handler) AddRocket(ctx context.Context, req *rkt.AddRocketRequest) (*rkt.AddRocketResponse, error) {
	log.Print("Add Rocket gRPC endpoint hit")

	if _, err := uuid.Parse(); err != nil {
		errorStatus := status.Error(codes.InvalidArgument, "uuid is not valid")
		log.Print("given uuid is not valid")

		return &rkt.AddRocketResponse{}, errorStatus
	}

	newRkt, err := h.RockerService.InsertRocket(ctx, rocket.Rocket{
		ID:   req.Rocket.Id,
		Type: req.Rocket.Type,
		Name: req.Rocket.Name,
	})
	if err != nil {
		log.Print("failed to insert rocket into database")
		return &rkt.AddRocketResponse{}, err
	}

	return &rkt.AddRocketResponse{
		Rocket: &rkt.Rocket{
			Id:   newRkt.ID,
			Type: newRkt.Type,
			Name: newRkt.Name,
		},
	}, nil
}

func (h Handler) DeleteRocket(ctx context.Context, req *rkt.DeleteRocketRequest) (*rkt.DeleteRocketResponse, error) {
	log.Print("delete rocket gRPC endpoint hit")

	err := h.RockerService.DeleteRocket(ctx, req.Id)
	if err != nil {
		return &rkt.DeleteRocketResponse{}, err
	}

	return &rkt.DeleteRocketResponse{Status: "successfully delete rocket"}, nil
}
