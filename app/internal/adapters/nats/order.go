package nats

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/shamank/wb-l0/app/internal/domain"
	"github.com/sirupsen/logrus"
)

func (n *Nats) createOrderFromMsg(msg *stan.Msg) {
	var order domain.Order

	if err := json.Unmarshal(msg.Data, &order); err != nil {
		logrus.Error(err.Error())
		return
	}

	if err := n.service.Orders.Create(n.ctx, order); err != nil {
		logrus.Error(err.Error())
		return
	}
}
