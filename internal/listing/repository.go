package listing

import (
	"context"
	"fmt"
	thumbnailspb "fuu/v/gen/go/grpc/thumbnails/v1"
	"fuu/v/internal/domain"
	"fuu/v/pkg/instrumentation"
	"time"

	"github.com/goccy/go-json"
	"github.com/marcopeocchi/fazzoletti/slices"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	rdb    *redis.Client
	logger *zap.SugaredLogger
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "listing.Count")
	defer span.End()

	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Directory{}).Count(&count).Error
	return count, err
}

func (r *Repository) Create(ctx context.Context, path, name, thumbnail string) (domain.Directory, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "listing.Create")
	defer span.End()

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
	_, span := otel.Tracer(otelName).Start(ctx, "listing.FindByName")
	defer span.End()

	m := domain.Directory{}
	err := r.db.WithContext(ctx).First(&m, name).Error
	return m, err
}

func (r *Repository) FindAllByName(ctx context.Context, filter string) (*[]domain.Directory, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "listing.FindAllByName")
	defer span.End()

	r.logger.Infow("FindAllByName", "filter", filter)
	all := new([]domain.Directory)

	cached, _ := r.rdb.Get(ctx, filter).Bytes()

	if len(cached) > 0 {
		json.Unmarshal(cached, all)
		instrumentation.CacheHitCounter.Add(1)
		return all, nil
	}

	err := r.db.WithContext(ctx).
		Joins("left join thumbnails on directories.path = thumbnails.folder").
		Where("name LIKE ?", "%"+filter+"%").
		Find(all).Error

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
	_, span := otel.Tracer(otelName).Start(ctx, "listing.FindAllRange")
	defer span.End()

	cacheKey := fmt.Sprint("findallrange", take, skip)

	r.logger.Infow("FindAllRange", "take", take, "skip", skip)
	_range := new([]domain.Directory)

	cached, _ := r.rdb.Get(ctx, cacheKey).Bytes()

	if len(cached) > 0 {
		json.Unmarshal(cached, _range)
		instrumentation.CacheHitCounter.Add(1)
		return _range, nil
	}

	conn, err := getGrpcClient("localhost:10099")
	if err != nil {
		span.End()
		r.logger.Fatalln(err)
		return nil, err
	}

	client := thumbnailspb.NewThumbnailServiceClient(conn)
	defer conn.Close()

	var _order string
	if order == domain.OrderByDate {
		_order = "updated_at desc"
	}
	if order == domain.OrderByName {
		_order = "name"
	}

	// err = r.db.WithContext(ctx).Order(_order).Limit(take).Offset(skip).Find(_range).Error
	err = r.db.WithContext(ctx).
		Joins("left join thumbnails on directories.path = thumbnails.folder").
		Order(_order).
		Limit(take).
		Offset(skip).
		Find(_range).Error

	if err != nil {
		return nil, err
	}

	paths := slices.Map(*_range, func(d domain.Directory) string {
		return d.Path
	})

	res, err := client.GetRange(ctx, &thumbnailspb.GetRangeRequest{
		Paths: paths,
	})

	for _, t := range res.Thumbnails {
		r.logger.Infoln(t.Path, t.Id)
	}

	if err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(*_range)
	if err != nil {
		return nil, err
	}

	err = r.rdb.SetNX(ctx, cacheKey, encoded, time.Minute).Err()
	r.logger.Warnw("FindAllRange", "warn", err)

	instrumentation.CacheMissCounter.Add(1)

	return _range, err
}

func (r *Repository) FindAll(ctx context.Context) (*[]domain.Directory, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "listing.FindAll")
	defer span.End()

	all := new([]domain.Directory)
	err := r.db.WithContext(ctx).Find(all).Error
	return all, err
}

func (r *Repository) Update(ctx context.Context, path, name, thumbnail string) (domain.Directory, error) {
	_, span := otel.Tracer(otelName).Start(ctx, "listing.Update")
	defer span.End()

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
	_, span := otel.Tracer(otelName).Start(ctx, "listing.Delete")
	defer span.End()

	m := domain.Directory{}
	err := r.db.WithContext(ctx).Where("path = ?", fmt.Sprintf("`%s`", path)).Delete(&domain.Directory{}).Error
	return m, err
}

// **** TESTING **** //

func getGrpcClient(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	}

	return grpc.DialContext(ctx, addr, opts...)
}
