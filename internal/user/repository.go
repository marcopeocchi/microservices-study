package user

import (
	"context"
	"errors"
	"fuu/v/internal/domain"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func (r *Repository) FindById(ctx context.Context, id uint) (domain.User, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "user.FindById")
	defer span.End()

	u := domain.User{}
	err := r.db.WithContext(ctx).First(&u, id).Error
	return u, err
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "user.FindByUsername")
	defer span.End()

	u := domain.User{}
	err := r.db.WithContext(ctx).First(&u, "username = ?", username).Error
	return u, err
}

func (r *Repository) Create(ctx context.Context, username, password string, role int) (domain.User, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "user.Create")
	defer span.End()

	if len(password) < 4 {
		return domain.User{}, errors.New("username must be at least 4 characters long")
	}
	if len(password) < 16 {
		return domain.User{}, errors.New("password must be at least 16 characters long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return domain.User{}, err
	}

	u := domain.User{
		Username: username,
		Password: string(hash),
		Role:     role,
	}

	err = r.db.WithContext(ctx).Create(&u).Error
	return u, err
}

func (r *Repository) Update(ctx context.Context, id uint, username, password string, role int) (domain.User, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "user.Update")
	defer span.End()

	u := domain.User{}
	err := r.db.WithContext(ctx).First(&u, id).Error

	if err != nil {
		return domain.User{}, err
	}

	u.Username = username
	u.Password = password
	u.Role = role

	err = r.db.WithContext(ctx).Save(&u).Error
	return u, err
}

func (r *Repository) Delete(ctx context.Context, id uint) (domain.User, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "user.Delete")
	defer span.End()

	u := domain.User{}
	err := r.db.WithContext(ctx).First(&u, id).Delete(&u, id).Error
	return u, err
}
