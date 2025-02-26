package route

import "context"

type Repository interface {
	Create(ctx context.Context, item *Route) error
	FindAll(ctx context.Context) ([]Route, error)
	FindOne(ctx context.Context, id int) (Route, error)
	Update(ctx context.Context, r Route) error
	Delete(ctx context.Context, id int) error
}
