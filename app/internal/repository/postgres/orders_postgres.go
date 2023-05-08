package postgres

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/shamank/wb-l0/app/internal/domain"
	"github.com/sirupsen/logrus"
)

type OrdersRepo struct {
	db *sqlx.DB
}

func NewOrdersRepo(db *sqlx.DB) *OrdersRepo {
	return &OrdersRepo{
		db: db,
	}
}

func (r *OrdersRepo) Create(ctx context.Context, order domain.Order) (domain.Order, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return domain.Order{}, err
	}

	orderQuery := `INSERT INTO "order" (
                     track_number,
                     entry, 
                     locale,
                     internal_signature,
                     customer_id, delivery_service,
                     shardkey,
                     sm_id,
                     oof_shard)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING order_uid`

	err = tx.QueryRow(
		orderQuery,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.OofShard,
	).Scan(&order.OrderUID)
	if err != nil {
		logrus.Error("error occurred in inserting order-table")
		tx.Rollback()
		return domain.Order{}, err
	}

	deliveryQuery := `INSERT INTO delivery
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.Exec(
		deliveryQuery,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		logrus.Error("error occurred in inserting delivery-table")
		tx.Rollback()
		return domain.Order{}, err
	}

	paymentQuery := `INSERT INTO payment
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(
		paymentQuery,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		logrus.Error("error occurred in inserting payment-table")
		tx.Rollback()
		return domain.Order{}, err
	}

	itemQuery := `INSERT INTO item(
						 order_uid,
						 track_number,
						 price,
						 rid,
						 name,
						 sale,
						 size,
						 total_price,
						 nm_id,
						 brand,
						 status)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	for _, item := range order.Items {
		_, err = tx.Exec(
			itemQuery,
			order.OrderUID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			logrus.Error("error occurred in inserting item-table")
			tx.Rollback()
			return domain.Order{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (r *OrdersRepo) Get(ctx context.Context, orderID string) (domain.Order, error) {

	var order domain.Order

	var items string

	query := `SELECT
             o.*,
             p.*,
             d.*,
             json_agg(
                     row_to_json(i
                         )
                 ) as items
         FROM
             "order" o
                 JOIN item i ON o.order_uid = i.order_uid
                 JOIN delivery d on o.order_uid = d.order_uid
                 JOIN payment p on o.order_uid = p.order_uid
         WHERE o.order_uid = $1
         GROUP BY
             o.order_uid, p.order_uid, d.order_uid`

	err := r.db.QueryRow(query, orderID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
		&order.OrderUID,
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDT,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
		&order.OrderUID,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
		&items,
	)
	if err != nil {
		return order, err
	}

	if err := json.Unmarshal([]byte(items), &order.Items); err != nil {
		return order, err
	}

	return order, nil
}

func (r *OrdersRepo) GetAll(ctx context.Context) ([]domain.Order, error) {

	orders := make([]domain.Order, 0)

	query := `SELECT
             o.*,
             p.*,
             d.*,
             json_agg(
                     row_to_json(i
                         )
                 ) as items
         FROM
             "order" o
                 JOIN item i ON o.order_uid = i.order_uid
                 JOIN delivery d on o.order_uid = d.order_uid
                 JOIN payment p on o.order_uid = p.order_uid
         GROUP BY
             o.order_uid, p.order_uid, d.order_uid`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order domain.Order
		var items string
		err = rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
			&order.OrderUID,
			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDT,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
			&order.OrderUID,
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,
			&items,
		)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(items), &order.Items); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
