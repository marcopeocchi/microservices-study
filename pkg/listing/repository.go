package listing

import (
	"context"
	"fmt"
	"fuu/v/pkg/domain"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	r.DB.WithContext(ctx).Model(&domain.Directory{}).Count(&count)
	return count, nil
}

func (r *Repository) Create(ctx context.Context, path, name, thumbnail string) (domain.Directory, error) {
	m := domain.Directory{
		Name:      name,
		Path:      path,
		Thumbnail: thumbnail,
		Loved:     false,
	}
	r.DB.WithContext(ctx).Create(&m)
	return m, nil
}

func (r *Repository) FindByName(ctx context.Context, name string) (domain.Directory, error) {
	m := domain.Directory{}
	r.DB.WithContext(ctx).First(&m, name)
	return m, nil
}

func (r *Repository) FindAllByName(ctx context.Context, filter string) (*[]domain.Directory, error) {
	all := new([]domain.Directory)
	r.DB.WithContext(ctx).Where("name LIKE ?", "%"+filter+"%").Find(all)
	return all, nil
}

func (r *Repository) FindAllRange(ctx context.Context, take, skip int) (*[]domain.Directory, error) {
	_range := new([]domain.Directory)
	r.DB.WithContext(ctx).Order("name").Limit(take).Offset(skip).Find(_range)
	return _range, nil
}

func (r *Repository) FindAll(ctx context.Context) (*[]domain.Directory, error) {
	all := new([]domain.Directory)
	r.DB.WithContext(ctx).Find(all)
	return all, nil
}

func (r *Repository) Update(ctx context.Context, path, name, thumbnail string) (domain.Directory, error) {
	m := domain.Directory{}
	r.DB.WithContext(ctx).First(&m)

	m.Name = name
	m.Path = path
	m.Thumbnail = thumbnail
	r.DB.WithContext(ctx).Save(&m)

	return m, nil
}

func (r *Repository) Delete(ctx context.Context, path string) (domain.Directory, error) {
	m := domain.Directory{}
	r.DB.WithContext(ctx).Where("path = ?", fmt.Sprintf("`%s`", path)).Delete(&domain.Directory{})
	return m, nil
}
