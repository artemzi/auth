package user

import (
	"github.com/artemzi/auth/internal/client/db"
	"github.com/artemzi/auth/internal/repository"
	"github.com/artemzi/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
