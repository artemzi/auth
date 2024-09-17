package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/artemzi/auth/internal/model"
	desc "github.com/artemzi/auth/pkg/user_v1"
)

func ToUserFromService(note *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        note.ID,
		Info:      ToUserInfoFromService(note.Info),
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromService(info model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            desc.Role(desc.Role_value[info.Role]),
	}
}

func ToUserInfoFromDesc(info *desc.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            info.Role.Enum().String(),
	}
}
