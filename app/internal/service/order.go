package service

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
	"github.com/shamank/wb-l0/app/internal/domain"
	"github.com/shamank/wb-l0/app/internal/repository"
	"github.com/sirupsen/logrus"
)

type OrdersService struct {
	memcache *cache.Cache
	repo     repository.Orders
	validate *validator.Validate
}

func NewOrdersService(repo repository.Orders, memcache *cache.Cache, validate *validator.Validate) *OrdersService {
	return &OrdersService{
		repo:     repo,
		memcache: memcache,
		validate: validate,
	}
}

func (s OrdersService) Create(ctx context.Context, order domain.Order) error {
	err := s.validate.Struct(order)
	if err != nil {
		logrus.Errorf("not valid order: %s", err.Error())
		return err
	}

	ord, err := s.repo.Create(ctx, order)
	if err != nil {
		return err
	}

	s.memcache.SetDefault(ord.OrderUID, ord)
	logrus.Infof("set %s to cache", ord.OrderUID)

	//
	//if err := s.repo.Create(ctx, order); err != nil {
	//	return err
	//}
	return nil
}

func (s OrdersService) Get(ctx context.Context, OrderID string) (domain.Order, error) {

	var order domain.Order

	data, ok := s.memcache.Get(OrderID)
	if !ok {
		order, err := s.repo.Get(ctx, OrderID)
		if err != nil {
			return domain.Order{}, err
		}

		s.memcache.SetDefault(order.OrderUID, order)
		logrus.Infof("set %s to cache", order.OrderUID)

		return order, nil

	}
	order, ok = data.(domain.Order)
	if !ok {
		return order, errors.New("problem with convert order from cache")
	}
	logrus.Infof("getting %s from cache", order.OrderUID)
	return order, nil
}

func (s OrdersService) GetAll(ctx context.Context) ([]domain.Order, error) {

	orders, err := s.repo.GetAll(ctx)
	if err != nil {
		logrus.Error("error with getall orders")
		return nil, err
	}
	logrus.Infof("getAll() in service")
	return orders, nil

}
