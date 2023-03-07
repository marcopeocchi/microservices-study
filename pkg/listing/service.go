package listing

import (
	"context"
	"fuu/v/pkg/domain"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	repo domain.ListingRepository
}

func (s *Service) ListAllDirectories(ctx context.Context) (*[]domain.DirectoryEnt, error) {
	ctx, span := trace.SpanFromContext(ctx).
		TracerProvider().
		Tracer("fs").
		Start(ctx, "listing.ListAllDirectories")

	defer span.End()

	dirs, err := s.repo.FindAll(ctx)
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

func (s *Service) ListAllDirectoriesRange(ctx context.Context, take, skip, order int) (*[]domain.DirectoryEnt, error) {
	ctx, span := trace.SpanFromContext(ctx).
		TracerProvider().
		Tracer("fs").
		Start(ctx, "listing.ListAllDirectoriesRange")

	defer span.End()

	dirs, err := s.repo.FindAllRange(ctx, take, skip, order)
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
			LastModified: dir.UpdatedAt,
		}
	}

	return &previews, nil
}

func (s *Service) ListAllDirectoriesLike(ctx context.Context, name string) (*[]domain.DirectoryEnt, error) {
	ctx, span := trace.SpanFromContext(ctx).
		TracerProvider().
		Tracer("fs").
		Start(ctx, "listing.ListAllDirectoriesLike")

	defer span.End()

	dirs, err := s.repo.FindAllByName(ctx, name)
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

func (s *Service) CountDirectories(ctx context.Context) (int64, error) {
	ctx, span := trace.SpanFromContext(ctx).
		TracerProvider().
		Tracer("fs").
		Start(ctx, "listing.CountDirectories")

	defer span.End()

	count, err := s.repo.Count(ctx)

	if err != nil {
		return 0, err
	}

	return count, nil
}
