package storage

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
)

type Repository[T any] struct {
	db *bun.DB
}

func NewRepository[T any](db *bun.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) ByID(ctx context.Context, id any) (*T, error) {
	model := new(T)
	err := r.db.NewSelect().Model(model).Where("id = ?", id).Scan(ctx)

	return model, err
}

func (r *Repository[T]) Create(ctx context.Context, model *T) (*T, error) {
	_, err := r.db.NewInsert().Model(model).Returning("*").Exec(ctx)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *Repository[T]) Update(ctx context.Context, model *T) (*T, error) {
	_, err := r.db.NewUpdate().Model(model).WherePK().Returning("*").Exec(ctx)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *Repository[T]) BulkCreate(ctx context.Context, models []T) error {
	if len(models) == 0 {
		return nil
	}

	_, err := r.db.NewInsert().Model(models).Exec(ctx)
	return err
}

func (r *Repository[T]) One(ctx context.Context, opts ...QueryOption) (*T, error) {
	model := new(T)

	cfg := &queryConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	q := r.db.NewSelect().Model(model)
	q = applyFilter(q, r.buildFilter(cfg))

	err := q.Limit(1).Scan(ctx)

	return model, err
}

func (r *Repository[T]) List(ctx context.Context, opts ...QueryOption) ([]*T, error) {
	var items []*T
	cfg := &queryConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	q := r.db.NewSelect().Model(items)
	q = applyFilter(q, r.buildFilter(cfg))
	if err := q.Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository[T]) Paginate(
	ctx context.Context,
	page int,
	perPage int,
	opts ...QueryOption,
) (*Paginator[*T], error) {
	page, perPage = NormalizePagination(page, perPage)

	var (
		items      []*T
		totalCount int
		totalPages int
	)
	g, cx := errgroup.WithContext(ctx)

	g.Go(func() error {
		cfg := &queryConfig{}
		for _, opt := range opts {
			opt(cfg)
		}

		cfg.limit = perPage
		cfg.offset = (page - 1) * perPage

		q := r.db.NewSelect().Model(&items)
		q = applyFilter(q, r.buildFilter(cfg))
		return q.Scan(cx)
	})

	g.Go(func() error {
		var err error
		cfg := &queryConfig{}
		for _, opt := range opts {
			opt(cfg)
		}
		cfg.limit = 1
		cfg.offset = 0
		cfg.orderBy = make([]Sort, 0)

		q := r.db.NewSelect().Model(&items)
		q = applyFilter(q, r.buildFilter(cfg))
		totalCount, err = q.Count(cx)

		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	totalPages = totalCount / perPage
	if totalCount%perPage != 0 {
		totalPages++
	}

	return &Paginator[*T]{
		Items:      items,
		Total:      totalCount,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}, nil
}

func applyFilter(q *bun.SelectQuery, f *Filter) *bun.SelectQuery {
	for _, r := range f.Rules {
		col := bun.Ident(r.Field)

		switch r.Operation {
		case OpEqual:
			q = q.Where("? = ?", col, r.Value)
		case OpNotEqual:
			q = q.Where("? != ?", col, r.Value)
		case OpGreaterThan:
			q = q.Where("? > ?", col, r.Value)
		case OpGreaterThanOrEqual:
			q = q.Where("? >= ?", col, r.Value)
		case OpLessThan:
			q = q.Where("? < ?", col, r.Value)
		case OpLessThanOrEqual:
			q = q.Where("? <= ?", col, r.Value)
		case OpLike:
			q = q.Where("? LIKE ?", col, r.Value)
		case OpILike:
			q = q.Where("? ILIKE ?", col, r.Value)
		case OpIn:
			q = q.Where("? IN (?)", col, bun.In(r.Value))
		case OpNotIn:
			q = q.Where("? NOT IN (?)", col, bun.In(r.Value))
		case OpIsNull:
			q = q.Where("? IS NULL", col)
		case OpIsNotNull:
			q = q.Where("? IS NOT NULL", col)
		case OpContains:
			q = q.Where("? @> ?", col, r.Value)
		case OpContainedBy:
			q = q.Where("? <@ ?", col, r.Value)
		case OpOverlaps:
			q = q.Where("? && ?", col, r.Value)
		case OpJsonContains:
			q = q.Where("? @> ?", col, r.Value)
		case OpJsonExists:
			q = q.Where("? ? ?", col, r.Value)
		}
	}

	for _, s := range f.OrderBy {
		dir := "ASC"
		if s.Descending {
			dir = "DESC"
		}
		q = q.OrderExpr(fmt.Sprintf("%s %s", s.Field, dir))
	}

	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	if f.Offset > 0 {
		q = q.Offset(f.Offset)
	}

	return q
}

func (r *Repository[T]) buildFilter(cfg *queryConfig) *Filter {
	if cfg == nil {
		return nil
	}

	return &Filter{
		Rules:   cfg.rules,
		OrderBy: cfg.orderBy,
		Limit:   cfg.limit,
		Offset:  cfg.offset,
	}
}
