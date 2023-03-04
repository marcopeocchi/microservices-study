package user

import (
	"context"
	"fuu/v/pkg/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func (r *Repository) FindById(ctx context.Context, id uint) (domain.User, error) {
	u := domain.User{}
	r.db.WithContext(ctx).First(&u, id)
	return u, nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	u := domain.User{}
	r.db.WithContext(ctx).First(&u, "username = ?", username)
	return u, nil
}

func (r *Repository) Create(ctx context.Context, username, password string, role int) (domain.User, error) {
	u := domain.User{
		Username: username,
		Password: password,
		Role:     role,
	}
	r.db.WithContext(ctx).Create(&u)
	return u, nil
}

func (r *Repository) Update(ctx context.Context, id uint, username, password string, role int) (domain.User, error) {
	u := domain.User{}
	r.db.WithContext(ctx).First(&u, id)

	u.Username = username
	u.Password = password
	u.Role = role

	r.db.WithContext(ctx).Save(&u)
	return u, nil
}

func (r *Repository) Delete(ctx context.Context, id uint) (domain.User, error) {
	u := domain.User{}
	r.db.WithContext(ctx).First(&u, id).Delete(&u, id)
	return u, nil
}
