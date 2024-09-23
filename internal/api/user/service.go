package user

import (
	"github.com/artemzi/auth/internal/service"
	desc "github.com/artemzi/auth/pkg/user_v1"
)

// TODO: add tests
// Implementation is an implementation of UserAPIV1Server.
type Implementation struct {
	desc.UnimplementedUserAPIV1Server
	userService service.UserService
}

// NewImplementation returns an implementation of UserAPIV1Server.
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
