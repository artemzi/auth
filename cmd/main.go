package main

import (
	"context"
	"flag"
	"net"
	"strconv"
	"time"

	"github.com/artemzi/auth/internal/config"
	desc "github.com/artemzi/auth/pkg/user_v1"
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
		log.Errorf("failed to INSERT user: %v", err)
		return nil, err
	}

	log.WithContext(ctx).Infof("inserted user with id: %d", userID)

	return &desc.CreateResponse{
		Id: userID,
	}, nil

}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	query := "SELECT id, name, email, role, created_at, updated_at FROM \"user\" WHERE id = $1;"

	var id int64
	var name, email, role string
	var createdAt, updatedAt time.Time

	err := s.pool.QueryRow(ctx, query, req.GetId()).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Errorf("failed to GET user: %v", err)
		return nil, err
	}

	log.WithContext(ctx).Info(color.GreenString("Got User id: "), req.GetId())

	roleVal, err := strconv.Atoi(role)
	if err != nil {
		log.Errorf("failed to parse role: %v", err)
		return nil, err
	}

	return &desc.GetResponse{
		User: &desc.User{
			Id: req.GetId(),
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  desc.Role(roleVal),
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		},
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	query := "UPDATE \"user\" SET name = $1, email = $2 WHERE id = $3;"

	_, err := s.pool.Exec(ctx, query, req.GetInfo().GetName().Value, req.GetInfo().GetEmail().Value, req.GetId())
	if err != nil {
		log.Errorf("failed to UPDATE user: %v", err)
		return nil, err
	}

	log.WithContext(ctx).Info(color.GreenString("Updated User id: "), req.GetId())
	out := new(emptypb.Empty)

	return out, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	query := "DELETE FROM \"user\" WHERE id = $1;"

	_, err := s.pool.Exec(ctx, query, req.GetId())
	if err != nil {
		log.Errorf("failed to GET user: %v", err)
		return nil, err
	}

	log.WithContext(ctx).Info(color.GreenString("Deleted User id: "), req.GetId())
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
