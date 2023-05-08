package nats

import (
	"context"
	"github.com/nats-io/stan.go"
	"github.com/shamank/wb-l0/app/config"
	"github.com/shamank/wb-l0/app/internal/service"
)

const (
	orderChannel = "channel-order"
)

type Nats struct {
	sc      stan.Conn
	service *service.Service
	ctx     context.Context
}

func NewNats(cfg config.NatsConfig, service *service.Service, ctx context.Context) (*Nats, error) {
	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsURL(cfg.URL))
	if err != nil {
		return nil, err
	}

	return &Nats{
		sc:      sc,
		service: service,
		ctx:     ctx,
	}, nil
}

func (n *Nats) InitSubscriptions() {
	n.sc.Subscribe(orderChannel, n.createOrderFromMsg, stan.DurableName("my-durable"))
}
