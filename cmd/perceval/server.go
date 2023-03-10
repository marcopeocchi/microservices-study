package main

import (
	"context"

	"fuu/v/cmd/perceval/model"
	thumbnailspb "fuu/v/gen/go/grpc/thumbnails/v1"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	otelName = "fuu/v/perceval/internal"
	format   = "webp"
)

type ThumbnailsService struct {
	db     *gorm.DB
	Logger *zap.SugaredLogger
}

func (t *ThumbnailsService) Generate(ctx context.Context, req *thumbnailspb.GenerateRequest) (*thumbnailspb.GenerateResponse, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "Generate")
	defer span.End()

	test := model.Thumbnail{}
	t.db.WithContext(ctx).First(&test, "folder = ?", req.Folder)

	if test.Thumbnail == "" {
		go convert(req.Path, req.Folder, req.Format, t.db, t.Logger)
	}

	return &thumbnailspb.GenerateResponse{
		Thumbnail: &thumbnailspb.Thumbnail{
			Id:     "",
			Path:   req.Path,
			Format: req.Format,
		},
	}, nil
}

func (t *ThumbnailsService) Delete(ctx context.Context, req *thumbnailspb.DeleteRequest) (*thumbnailspb.DeleteResponse, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "Delete")
	defer span.End()

	//TODO: implementazione
	err := delete(req.Path, t.Logger)
	if err != nil {
		return nil, err
	}

	return &thumbnailspb.DeleteResponse{
		Thumbnail: &thumbnailspb.Thumbnail{
			Id:     "",
			Path:   req.Path,
			Format: format,
		},
	}, nil
}

func (t *ThumbnailsService) Get(ctx context.Context, req *thumbnailspb.GetRequest) (*thumbnailspb.GetResponse, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "Get")
	defer span.End()

	//TODO: implementazione
	id, path, err := getByPath(req.Path, t.Logger)
	if err != nil {
		return nil, err
	}

	return &thumbnailspb.GetResponse{
		Thumbnail: &thumbnailspb.Thumbnail{
			Id:     id,
			Path:   path,
			Format: format,
		},
	}, nil
}

func (t *ThumbnailsService) GetRange(ctx context.Context, req *thumbnailspb.GetRangeRequest) (*thumbnailspb.GetRangeResponse, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "Get")
	defer span.End()

	//TODO: implementazione
	ids, err := getManyByPath(ctx, req.Paths, t.db, t.Logger)
	if err != nil {
		return nil, err
	}

	res := make([]*thumbnailspb.Thumbnail, len(*ids))

	for i, pair := range *ids {
		res[i] = &thumbnailspb.Thumbnail{
			Id:   pair.Thumbnail,
			Path: pair.Path,
		}
	}

	return &thumbnailspb.GetRangeResponse{
		Thumbnails: res,
	}, nil
}
