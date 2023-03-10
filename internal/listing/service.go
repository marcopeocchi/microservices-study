package listing

import (
	"context"
	"fuu/v/internal/domain"

	"go.opentelemetry.io/otel"
)

const otelName = "fuu/v/internal/listing"

type Service struct {
	repo domain.ListingRepository
}

func (s *Service) ListAllDirectories(ctx context.Context) (*[]domain.DirectoryEnt, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "listing.ListAllDirectories")
	defer span.End()

	dirs, err := s.repo.FindAll(ctx)
	if err != nil {
		span.RecordError(err)
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
	_, span := otel.Tracer(otelName).Start(ctx, "listing.ListAllDirectoriesRange")
	defer span.End()

	dirs, err := s.repo.FindAllRange(ctx, take, skip, order)
	if err != nil {
		span.RecordError(err)
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

func (s *Service) ListAllDirectoriesLike(ctx context.Context, name string, take, skip int) (*[]domain.DirectoryEnt, int64, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "listing.ListAllDirectoriesLike")
	defer span.End()

	if len(name) <= 2 {
		return &[]domain.DirectoryEnt{}, 0, nil
	}

	dirs, count, err := s.repo.FindLikeNameRange(ctx, name, take, skip)
	if err != nil {
		span.RecordError(err)
		return nil, 0, err
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

	return &previews, count, nil
}

func (s *Service) CountDirectories(ctx context.Context) (int64, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "listing.CountDirectories")
	defer span.End()

	count, err := s.repo.Count(ctx)

	if err != nil {
		span.RecordError(err)
		return 0, err
	}

	return count, nil
}
