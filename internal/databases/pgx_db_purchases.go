package databases

import (
	"balance/internal/models"

	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

// Purchase performs the purchase transaction for given orderId, userId, serviceId and amount.
//
// 1) checks if money were reserved or purchase has already happened
//
// 2) if everything is ok sets reserve status to purchased
//
// 3) writes report to operations table
func (p PgxDB) Purchase(userId, serviceId, orderId uint64, amount int64) error {
	ctx := context.Background()

	var err error
	defer func() {
		if err != nil {
			p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: purchase: %v", err), nil)
		}
	}()

	// start transaction and defer its closing
	tx, err := p.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	reserve := models.Reserve{
		UserID:    userId,
		ServiceID: serviceId,
		OrderID:   orderId,
	}

	// find reserve
	err = tx.QueryRow(ctx, "select * from reserves where order_id = $1 and user_id = $2 and service_id = $3",
		reserve.OrderID, reserve.UserID, reserve.ServiceID).Scan(&reserve.OrderID, &reserve.UserID, &reserve.ServiceID,
		&reserve.Amount, &reserve.Purchased, &reserve.ReservedAt, &reserve.PurchasedAt)
	if err != nil && errors.Is(err, pgx.ErrNoRows) { // reserve not found
		err = fmt.Errorf("db: get purchase: money were not reserved for order %d, user %d and service %d", userId, serviceId, orderId)
		return err
	} else if err != nil {
		return err
	} else if reserve.Amount != amount { // wrong amount
		got := (float32(amount%100) * 0.01) + float32(amount/100)
		stored := (float32(reserve.Amount%100) * 0.01) + float32(reserve.Amount/100)
		err = fmt.Errorf("db: purchase: wrong purchase amount, stored in reserve: %.2f, got: %.2f", stored, got)
		return err
	} else if reserve.Purchased { // already purchased
		err = errors.New("db: purchase: the purchase has already happened")
		return err
	}

	loc, _ := time.LoadLocation("Europe/Moscow")
	purchasedAt := time.Now().In(loc)
	reserve.PurchasedAt = &purchasedAt
	reserve.Purchased = true

	// update purchase status
	var updateId uint64
	err = tx.QueryRow(ctx, "update reserves SET purchased = $1, purchased_at = $2 where order_id = $3 returning order_id",
		reserve.Purchased, reserve.PurchasedAt, reserve.OrderID).Scan(&updateId)
	if err != nil {
		return err
	}

	// write to operations table
	var service models.Service
	err = tx.QueryRow(ctx, "select * from services where id = $1;", serviceId).Scan(&service.ID, &service.Name)
	if err != nil {
		return err
	}
	reserve.Service = service

	var createdId uint64
	err = tx.QueryRow(ctx, "insert into operations (user_id, service_id, service_name, amount, done_at) values ($1, $2, $3, $4, $5) returning id",
		reserve.UserID, reserve.ServiceID, reserve.Service.Name, reserve.Amount*-1, purchasedAt).Scan(&createdId)
	if err != nil {
		return err
	}

	return err
}
