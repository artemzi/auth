package user

import (
	"github.com/artemzi/auth/internal/service"
	desc "github.com/artemzi/auth/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserAPIV1Server
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
