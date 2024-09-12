package main

import (
	"context"
	"flag"
	"net"

	"github.com/artemzi/auth/internal/config"
	desc "github.com/artemzi/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserAPIV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	query := "INSERT INTO \"user\" (name, email, password, password_confirm, role) VALUES ($1, $2, $3, $4, $5) RETURNING id;"

	var userID int64
	err := s.pool.QueryRow(
		ctx,
		query,
		req.GetInfo().Name, req.GetInfo().Email, req.GetInfo().Password, req.GetInfo().PasswordConfirm, req.GetInfo().Role).
		Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.WithContext(ctx).Infof("inserted user with id: %d", userID)

	return &desc.CreateResponse{
		Id: userID,
	}, nil

}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.WithContext(ctx).Info(color.GreenString("Get User id: "), req.GetId())

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
	log.WithContext(ctx).Info(color.GreenString("Update User id: "), req.GetId())
	out := new(emptypb.Empty)

	return out, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.WithContext(ctx).Info(color.GreenString("Delete User id: "), req.GetId())
	out := new(emptypb.Empty)

	return out, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserAPIV1Server(s, &server{pool: pool})

	log.Info(color.GreenString("server listening at "), lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
