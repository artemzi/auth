package mapper

import (
	"github.com/artemzi/auth/internal/model"
	entity "github.com/artemzi/auth/internal/repository/user/entity"
)

func ToUserFromEntity(user *entity.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromEntity(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserInfoFromEntity(info entity.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Passowrd,
		PasswordConfirm: info.PasswordConfirm,
		Role:            info.Role,
	}
}
