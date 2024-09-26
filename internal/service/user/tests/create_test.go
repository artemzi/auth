package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/artemzi/auth/internal/model"
	"github.com/artemzi/auth/internal/repository"
	repositoryMock "github.com/artemzi/auth/internal/repository/mocks"
	"github.com/artemzi/auth/internal/service/user"
	desc "github.com/artemzi/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userRepositoryMokj func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx  context.Context
		info *model.UserInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id              = gofakeit.Int64()
		name            = gofakeit.Name()
		email           = gofakeit.Email()
		password        = "secret"
		passwordConfirm = "secret"
		role            = desc.Role_ROLE_USER

		serviceErr = fmt.Errorf("service error")

		req = &model.UserInfo{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            role.String(),
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMokj
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				info: req,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMock.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx:  ctx,
				info: req,
			},
			want: 0,
			err:  serviceErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMock.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.userRepositoryMock(mc)
			service := user.NewMockService(userRepositoryMock)

			result, err := service.Create(tt.args.ctx, tt.args.info)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, result)
		})
	}
}
