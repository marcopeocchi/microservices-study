package user

import (
	"context"
	"errors"
	"fuu/v/pkg/domain"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   domain.UserRepository
	logger *zap.SugaredLogger
}

func (s *Service) Login(ctx context.Context, username, password string) (*string, error) {
	ctx, span := trace.SpanFromContext(ctx).
		TracerProvider().
		Tracer("fs").
		Start(ctx, "user.Login")

	defer span.End()

	u, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if u.Username == "" {
		return nil, errors.New("username not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, errors.New("password bcrypt hash does not match")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    u.ID,
		"username":  u.Username,
		"role":      u.Role,
		"expiresAt": time.Now().Add(time.Minute * 30),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWTSECRET")))

	return &tokenString, err
}

func (s *Service) Create(ctx context.Context, username, password string, role int) (domain.User, error) {
	ctx, span := trace.SpanFromContext(ctx).
		TracerProvider().
		Tracer("fs").
		Start(ctx, "user.Create")

	defer span.End()

	u, err := s.repo.Create(ctx, username, password, role)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
