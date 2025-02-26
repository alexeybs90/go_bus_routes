package station

import "context"

type Repository interface {
	Create(ctx context.Context, item *Station) error
	FindAll(ctx context.Context) ([]Station, error)
	FindOne(ctx context.Context, id int) (Station, error)
	Update(ctx context.Context, r Station) error
	Delete(ctx context.Context, id int) error
}
