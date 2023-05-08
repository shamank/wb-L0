package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/shamank/wb-l0/app/internal/domain"
	"github.com/shamank/wb-l0/app/internal/repository/postgres"
)

type Orders interface {
	Create(ctx context.Context, order domain.Order) (domain.Order, error)
	Get(ctx context.Context, orderID string) (domain.Order, error)
	GetAll(ctx context.Context) ([]domain.Order, error)
}

type Repository struct {
	Orders Orders
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Orders: postgres.NewOrdersRepo(db),
	}
}
