package main

import (
	"context"
	"fmt"
	"net"

	desc "github.com/artemzi/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserAPIV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Info(color.GreenString("Create User email: "), req.GetInfo().Email)

	return &desc.CreateResponse{Id: int64(gofakeit.Uint64())}, nil

}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Info(color.GreenString("Get User id: "), req.GetId())

	return &desc.GetResponse{
		User: &desc.User{
			Id: req.GetId(),
			Info: &desc.UserInfo{
				Name:            gofakeit.BeerName(),
				Email:           gofakeit.Email(),
				Password:        "123",
				PasswordConfirm: "123",
				Role:            desc.Role_ROLE_ADMIN,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Info(color.GreenString("Update User id: "), req.GetId())
	out := new(emptypb.Empty)

	return out, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Info(color.GreenString("Delete User id: "), req.GetId())
	out := new(emptypb.Empty)

	return out, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserAPIV1Server(s, &server{})

	log.Info(color.GreenString("server listening at "), lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
