//go:build wireinject
// +build wireinject

package listing

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func Wire(db *gorm.DB) *Handler {
	panic(wire.Build(ProviderSet))
}
