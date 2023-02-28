package listing

import (
	"context"
	"fmt"
	"fuu/v/pkg/domain"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	r.db.WithContext(ctx).Model(&domain.Directory{}).Count(&count)
	return count, nil
}

func (r *Repository) Create(ctx context.Context, path, name, thumbnail string) (domain.Directory, error) {
	m := domain.Directory{
		Name:      name,
		Path:      path,
		Thumbnail: thumbnail,
		Loved:     false,
	}
	r.db.WithContext(ctx).Create(&m)
	return m, nil
}

func (r *Repository) FindByName(ctx context.Context, name string) (domain.Directory, error) {
	m := domain.Directory{}
	r.db.WithContext(ctx).First(&m, name)
	return m, nil
}

func (r *Repository) FindAllByName(ctx context.Context, filter string) (*[]domain.Directory, error) {
	all := new([]domain.Directory)

	cached, _ := r.rdb.Get(ctx, filter).Bytes()
	if len(cached) > 0 {
		json.Unmarshal(cached, all)
		return all, nil
	}

	r.db.WithContext(ctx).Where("name LIKE ?", "%"+filter+"%").Find(all)

	encoded, err := json.Marshal(*all)
	if err != nil {
		return nil, err
	}
	r.rdb.SetNX(ctx, filter, encoded, time.Minute)

	return all, nil
}

func (r *Repository) FindAllRange(ctx context.Context, take, skip, order int) (*[]domain.Directory, error) {
	_range := new([]domain.Directory)

	var _order string
	if order == domain.OrderByDate {
		_order = "updated_at desc"
	}
	if order == domain.OrderByName {
		_order = "name"
	}

	r.db.WithContext(ctx).Order(_order).Limit(take).Offset(skip).Find(_range)
	return _range, nil
}

func (r *Repository) FindAll(ctx context.Context) (*[]domain.Directory, error) {
	all := new([]domain.Directory)
	r.db.WithContext(ctx).Find(all)
	return all, nil
}

func (r *Repository) Update(ctx context.Context, path, name, thumbnail string) (domain.Directory, error) {
	m := domain.Directory{}
	r.db.WithContext(ctx).First(&m)

	m.Name = name
	m.Path = path
	m.Thumbnail = thumbnail
	r.db.WithContext(ctx).Save(&m)

	return m, nil
}

func (r *Repository) Delete(ctx context.Context, path string) (domain.Directory, error) {
	m := domain.Directory{}
	r.db.WithContext(ctx).Where("path = ?", fmt.Sprintf("`%s`", path)).Delete(&domain.Directory{})
	return m, nil
}
