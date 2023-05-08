package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
	"github.com/shamank/wb-l0/app/internal/domain"
	"github.com/shamank/wb-l0/app/internal/repository"
)

type Orders interface {
	Create(ctx context.Context, order domain.Order) error
	Get(ctx context.Context, OrderID string) (domain.Order, error)
	GetAll(ctx context.Context) ([]domain.Order, error)
}

type Service struct {
	Orders Orders
}

func NewService(repo *repository.Repository, memcache *cache.Cache) *Service {
	validate := validator.New()
	return &Service{
		Orders: NewOrdersService(repo.Orders, memcache, validate),
	}
}
