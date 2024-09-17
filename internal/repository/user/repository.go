package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/artemzi/auth/internal/client/db"
	"github.com/artemzi/auth/internal/model"
	"github.com/artemzi/auth/internal/repository"
	entity "github.com/artemzi/auth/internal/repository/user/entity"
	"github.com/artemzi/auth/internal/repository/user/mapper"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

const (
	tableName = "user"

	idColumn              = "id"
	nameColumn            = "name"
	emailColumn           = "email"
	passwordColumn        = "password"
	passwordConfirmColumn = "password_confirm"
	roleColumn            = "role"
	createdAtColumn       = "created_at"
	updatedAtColumn       = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, passwordConfirmColumn, roleColumn).
		Values(info.Name, info.Email, info.Password, info.PasswordConfirm, info.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var userID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
	if err != nil {
		log.Errorf("failed to INSERT user: %v", err)
		return 0, err
	}

	log.WithContext(ctx).Infof("inserted user with id: %d", userID)
	return userID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	query := "SELECT id, name, email, role, created_at, updated_at FROM \"user\" WHERE id = $1;"

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user entity.User
	err := r.db.DB().QueryRowContext(ctx, q, id).
		Scan(&user.ID, &user.Info.Name, &user.Info.Email, &user.Info.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Errorf("failed to GET user: %v", err)
		return nil, err
	}

	log.WithContext(ctx).Info(color.GreenString("Got User id: "), id)
	return mapper.ToUserFromEntity(&user), nil
}

// func (s *repo) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
// 	query := "UPDATE \"user\" SET name = $1, email = $2 WHERE id = $3;"

// 	_, err := s.pool.Exec(ctx, query, req.GetInfo().GetName().Value, req.GetInfo().GetEmail().Value, req.GetId())
// 	if err != nil {
// 		log.Errorf("failed to UPDATE user: %v", err)
// 		return nil, err
// 	}

// 	log.WithContext(ctx).Info(color.GreenString("Updated User id: "), req.GetId())
// 	out := new(emptypb.Empty)

// 	return out, nil
// }

// func (r *repo) Delete(ctx context.Context, id int64) error {
// 	query := "DELETE FROM \"user\" WHERE id = $1;"

// 	_, err := r.db.DB().pool.Exec(ctx, query, id)
// 	if err != nil {
// 		log.Errorf("failed to DELETE user: %v", err)
// 		return err
// 	}

// 	log.WithContext(ctx).Info(color.GreenString("Deleted User id: "), id)
// 	return nil
// }
