package user

import (
	"context"
	"log"

	"github.com/artemzi/auth/internal/mapper"
	desc "github.com/artemzi/auth/pkg/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, mapper.ToUserInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
