package listing

import (
	"context"
	"fmt"
	"fuu/v/internal/domain"
	"fuu/v/pkg/instrumentation"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	rdb    *redis.Client
	logger *zap.SugaredLogger
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Directory{}).Count(&count).Error
	return count, err
}

func (r *Repository) Create(ctx context.Context, path, name, thumbnail string) (domain.Directory, error) {
	m := domain.Directory{
		Name:      name,
		Path:      path,
		Thumbnail: thumbnail,
		Loved:     false,
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	return m, err
}

func (r *Repository) FindByName(ctx context.Context, name string) (domain.Directory, error) {
	m := domain.Directory{}
	err := r.db.WithContext(ctx).First(&m, name).Error
	return m, err
}

func (r *Repository) FindAllByName(ctx context.Context, filter string) (*[]domain.Directory, error) {
	r.logger.Infow("FindAllByName", "filter", filter)
	all := new([]domain.Directory)

	cached, _ := r.rdb.Get(ctx, filter).Bytes()

	if len(cached) > 0 {
		json.Unmarshal(cached, all)
		instrumentation.CacheHitCounter.Add(1)
		return all, nil
	}

	err := r.db.WithContext(ctx).Where("name LIKE ?", "%"+filter+"%").Find(all).Error
	if err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(*all)
	if err != nil {
		return nil, err
	}
	err = r.rdb.SetNX(ctx, filter, encoded, time.Minute).Err()
	r.logger.Warnw("FindAllRange", "warn", err)

	instrumentation.CacheMissCounter.Add(1)

	return all, nil
}

func (r *Repository) FindAllRange(ctx context.Context, take, skip, order int) (*[]domain.Directory, error) {
	r.logger.Infow("FindAllRange", "take", take, "skip", skip)
	_range := new([]domain.Directory)

	var _order string
	if order == domain.OrderByDate {
		_order = "updated_at desc"
	}
	if order == domain.OrderByName {
		_order = "name"
	}

	err := r.db.WithContext(ctx).Order(_order).Limit(take).Offset(skip).Find(_range).Error
	return _range, err
}

func (r *Repository) FindAll(ctx context.Context) (*[]domain.Directory, error) {
	all := new([]domain.Directory)
	err := r.db.WithContext(ctx).Find(all).Error
	return all, err
}

func (r *Repository) Update(ctx context.Context, path, name, thumbnail string) (domain.Directory, error) {
	m := domain.Directory{}
	err := r.db.WithContext(ctx).First(&m).Error
	if err != nil {
		return domain.Directory{}, err
	}

	m.Name = name
	m.Path = path
	m.Thumbnail = thumbnail
	err = r.db.WithContext(ctx).Save(&m).Error

	return m, err
}

func (r *Repository) Delete(ctx context.Context, path string) (domain.Directory, error) {
	m := domain.Directory{}
	err := r.db.WithContext(ctx).Where("path = ?", fmt.Sprintf("`%s`", path)).Delete(&domain.Directory{}).Error
	return m, err
}
