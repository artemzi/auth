package user

import (
	"github.com/artemzi/auth/internal/client/db"
	"github.com/artemzi/auth/internal/repository"
	def "github.com/artemzi/auth/internal/service"
)

var _ def.UserService = (*serv)(nil)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager) *serv {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
