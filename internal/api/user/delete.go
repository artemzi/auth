package user

import (
	"context"
	"log"

	desc "github.com/artemzi/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete implements user.UserServer
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("deleted user with id: %d", req.GetId())

	out := new(emptypb.Empty)
	return out, nil
}
