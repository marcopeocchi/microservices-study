package listing

import (
	"context"
	"fuu/v/pkg/domain"
)

var (
	ctx = context.Background()
)

type Service struct {
	Repo domain.ListingRepository
}

func (s *Service) ListAllDirectories() (*[]domain.DirectoryEnt, error) {
	dirs, err := s.Repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	previews := make([]domain.DirectoryEnt, len(*dirs))

	for i, dir := range *dirs {
		previews[i] = domain.DirectoryEnt{
			Id:        dir.ID,
			Name:      dir.Name,
			Loved:     dir.Loved,
			Thumbnail: dir.Thumbnail,
		}
	}

	return &previews, nil
}

func (s *Service) ListAllDirectoriesRange(take, skip int) (*[]domain.DirectoryEnt, error) {
	dirs, err := s.Repo.FindAllRange(ctx, take, skip)
	if err != nil {
		return nil, err
	}

	previews := make([]domain.DirectoryEnt, len(*dirs))

	for i, dir := range *dirs {
		previews[i] = domain.DirectoryEnt{
			Id:        dir.ID,
			Name:      dir.Name,
			Loved:     dir.Loved,
			Thumbnail: dir.Thumbnail,
		}
	}

	return &previews, nil
}

func (s *Service) ListAllDirectoriesLike(name string) (*[]domain.DirectoryEnt, error) {
	dirs, err := s.Repo.FindAllByName(ctx, name)
	if err != nil {
		return nil, err
	}

	previews := make([]domain.DirectoryEnt, len(*dirs))

	for i, dir := range *dirs {
		previews[i] = domain.DirectoryEnt{
			Id:           dir.ID,
			Name:         dir.Name,
			Loved:        dir.Loved,
			Thumbnail:    dir.Thumbnail,
			LastModified: dir.CreatedAt,
		}
	}

	return &previews, nil
}

func (s *Service) CountDirectories() (int64, error) {
	count, err := s.Repo.Count(ctx)

	if err != nil {
		return 0, err
	}

	return count, nil
}
