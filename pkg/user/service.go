package user

import (
	"context"
	"errors"
	"fuu/v/pkg/common"
	"fuu/v/pkg/domain"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo domain.UserRepository
}

func (s *Service) Login(ctx context.Context, username, password string) (*string, error) {
	u, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if u.Username == "" {
		return nil, errors.New("username not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    u.ID,
		"username":  u.Username,
		"role":      u.Role,
		"expiresAt": common.TOKEN_EXPIRE_TIME,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWTSECRET")))

	return &tokenString, err
}

func (s *Service) Create(ctx context.Context, username, password string, role int) (domain.User, error) {
	u, err := s.repo.Create(ctx, username, password, role)
	if err != nil {
		return domain.User{}, nil
	}
	return u, nil
}
