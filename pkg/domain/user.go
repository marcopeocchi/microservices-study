package domain

import (
	"context"
	"net/http"

	"gorm.io/gorm"
)

type UserRole int

const (
	Standard UserRole = iota
	Admin
)

type UserEnt struct {
	ID       uint
	Username string
	Password string
	Role     int
}

type User struct {
	gorm.Model
	ID       uint
	Username string
	Password string
	Role     int
}

type UserRepository interface {
	FindById(ctx context.Context, id uint) (User, error)

	FindManyById(ctx context.Context, id uint) (*[]User, error)

	FindByUsername(ctx context.Context, username string) (User, error)

	Create(ctx context.Context, username, password string, role int) (User, error)

	Update(ctx context.Context, username, password string, role int) (User, error)

	Delete(ctx context.Context, id uint) (User, error)
}

type UserService interface {
	FindById(id uint) (User, error)

	FindByUsername(username string) (User, error)
}

type UserHandler interface {
	Login() http.HandlerFunc

	Logout() http.HandlerFunc
}
