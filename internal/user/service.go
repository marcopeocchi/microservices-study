package user

import (
	"context"
	"errors"
	"fuu/v/internal/domain"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const otelName = "fuu/v/internal/user"

type Service struct {
	repo   domain.UserRepository
	logger *zap.SugaredLogger
}

func (s *Service) Login(ctx context.Context, username, password string) (*string, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "user.Login")
	defer span.End()

	u, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	if u.Username == "" {
		err := errors.New("username not found")
		span.RecordError(err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		span.RecordError(err)
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
	_, span := otel.Tracer(otelName).Start(ctx, "user.Create")
	defer span.End()

	u, err := s.repo.Create(ctx, username, password, role)
	if err != nil {
		span.RecordError(err)
		return domain.User{}, err
	}
	return u, nil
}
