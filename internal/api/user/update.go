package user

import (
	"context"
	"log"

	"github.com/artemzi/auth/internal/mapper"
	desc "github.com/artemzi/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Update updates user
func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, mapper.ToUserFromUpdateDesc(req))
	if err != nil {
		return nil, err
	}

	log.Printf("updated user with id: %d", req.GetId())

	out := new(emptypb.Empty)
	return out, nil
}
