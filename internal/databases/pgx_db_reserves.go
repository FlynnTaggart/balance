package databases

import (
	"balance/internal/models"

	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

// Reserve performs money reserve transaction for given orderId, userId, serviceId and amount.
//
// 1) checks if given user and service exist
//
// 2) subtracts user balance by amount
//
// 3) writes into reserves table with purchased status = false
func (p PgxDB) Reserve(userId, serviceId, orderId uint64, amount int64) error {
	ctx := context.Background()

	var err error
	defer func() {
		if err != nil {
			p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: reserve: %v", err), nil)
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

	// check user by id
	var user models.User
	err = tx.QueryRow(ctx, "select * from users where id = $1;", userId).Scan(&user.ID, &user.Balance)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("db: reserve: no such user with id %d", userId)
		return err
	} else if err != nil {
		return err
	} else if amount > user.Balance {
		needed := (float32(amount%100) * 0.01) + float32(amount/100)
		got := (float32(user.Balance%100) * 0.01) + float32(user.Balance/100)
		err = fmt.Errorf("db: reserve: the user %d doesn't have enough money, needed: %.2f, user has: %.2f", userId, needed, got)
		return err
	}

	// check service by id
	var checkServiceId uint64
	err = tx.QueryRow(ctx, "select id from services where id = $1;", serviceId).Scan(&checkServiceId)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("db: reserve: no such service with id %d", serviceId)
		return err
	} else if err != nil {
		return err
	}

	// subtract user's balance and update user
	user.Balance -= amount
	var updateId uint64
	err = tx.QueryRow(ctx, "update users SET balance = $2 where id = $1 returning id", user.ID, user.Balance).Scan(&updateId)
	if err != nil {
		return err
	}

	// set time location
	loc, _ := time.LoadLocation("Europe/Moscow")
	date := time.Now().In(loc)

	// insert into reserves table
	var reserveId uint64
	err = tx.QueryRow(ctx, "insert into reserves (order_id, user_id, service_id, amount, purchased, reserved_at) values ($1, $2, $3, $4, $5, $6) returning order_id",
		orderId, userId, serviceId, amount, false, date).Scan(&reserveId)
	if err != nil {
		return err
	}

	// start a new goroutine that returns money to user if the order was not purchased with a time (10 minutes by default)
	// TODO: make delete reserve timeout configurable
	go func() {
		time.Sleep(10 * time.Minute)
		err = p.DeleteReserve(userId, serviceId, orderId, amount)
		p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: reserve: %v", err), nil)
	}()

	return err
}

// GetReserve returns reserve by given userId, serviceId, orderId
func (p PgxDB) GetReserve(userId, serviceId, orderId uint64) (models.Reserve, error) {
	ctx := context.Background()

	var err error
	defer func() {
		if err != nil {
			p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: get reserve: %v", err), nil)
		}
	}()

	// start read only transaction and defer its closing
	tx, err := p.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		DeferrableMode: pgx.Deferrable,
		AccessMode:     pgx.ReadOnly,
	})
	if err != nil {
		return models.Reserve{}, err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	reserve := models.Reserve{
		OrderID:   orderId,
		UserID:    userId,
		ServiceID: serviceId,
	}

	// it would be better to select only amount and reservedAt but pgx gets error then
	err = tx.QueryRow(ctx, "select * from reserves where order_id = $1",
		reserve.OrderID).Scan(&reserve.OrderID, &reserve.UserID, &reserve.ServiceID,
		&reserve.Amount, &reserve.Purchased, &reserve.ReservedAt, &reserve.PurchasedAt)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("db: get reserve: money were not reserved for order %d", orderId)
		return models.Reserve{}, err
	} else if err != nil {
		return models.Reserve{}, err
	}

	// get user struct
	var user models.User
	err = tx.QueryRow(ctx, "select * from users where id = $1;", userId).Scan(&user.ID, &user.Balance)
	if err != nil {
		return models.Reserve{}, err
	}
	reserve.User = user

	// get service struct
	var service models.Service
	err = tx.QueryRow(ctx, "select * from services where id = $1;", serviceId).Scan(&service.ID, &service.Name)
	if err != nil {
		return models.Reserve{}, err
	}
	reserve.Service = service

	return reserve, err
}

// DeleteReserve deletes reserve by given userId, serviceId, orderId and amount
//
// returns reserved money to user
func (p PgxDB) DeleteReserve(userId, serviceId, orderId uint64, amount int64) error {
	ctx := context.Background()

	var err error
	defer func() {
		if err != nil {
			p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: delete reserve: %v", err), nil)
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

	// get the reserve
	err = tx.QueryRow(ctx, "select * from reserves where order_id = $1 and user_id = $2 and service_id = $3",
		reserve.OrderID, reserve.UserID, reserve.ServiceID).Scan(&reserve.OrderID, &reserve.UserID, &reserve.ServiceID,
		&reserve.Amount, &reserve.Purchased, &reserve.ReservedAt, &reserve.PurchasedAt)
	if err != nil {
		return err
	}
	// if reserve found, and it is uncompleted then return money and delete the found row
	if !reserve.Purchased {
		// get user
		var user models.User
		err = tx.QueryRow(ctx, "select * from users where id = $1;", userId).Scan(&user.ID, &user.Balance)
		if err != nil {
			return err
		}

		// return money
		user.Balance += amount

		// update user
		var updateId uint64
		err = tx.QueryRow(ctx, "update users SET balance = $2 where id = $1 returning id", user.ID, user.Balance).Scan(&updateId)
		if err != nil {
			return err
		}

		// delete reserve
		var deleteId uint64
		err = tx.QueryRow(ctx, "delete from reserves where order_id = $1 returning order_id",
			orderId).Scan(&deleteId)
		if err != nil {
			return err
		}
	}
	return err
}
