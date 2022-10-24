package databases

import (
	"balance/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
	"strings"
	"time"
)

type PgxDB struct {
	*pgxpool.Pool
}

func NewPgxDB(pool *pgxpool.Pool) *PgxDB {
	return &PgxDB{Pool: pool}
}

func (p PgxDB) GetBalance(id uint64) (float32, error) {
	ctx := context.TODO()

	var balance float32
	err := p.QueryRow(ctx, "select balance from users where id = $1;", id).Scan(&balance)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("db: reserve: no such user with id %d", id)
	} else if err != nil {
		return 0, err
	}
	return balance, err
}

func (p PgxDB) AddBalance(id uint64, amount float32) error {
	ctx := context.TODO()

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

	var user models.User
	err = tx.QueryRow(ctx, "select * from users where id = $1;", id).Scan(&user.ID, &user.Balance)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		user = models.User{
			ID:      id,
			Balance: amount,
		}
		var createdId uint64
		err = tx.QueryRow(ctx, "insert into users (id, balance) values ($1, $2) returning id", user.ID, user.Balance).Scan(&createdId)
		if err != nil {
			return err
		}
		return err
	} else if err != nil {
		return err
	}

	user.Balance += amount

	var updateId uint64
	err = tx.QueryRow(ctx, "update users SET balance = $2 where id = $1 returning id", user.ID, user.Balance).Scan(&updateId)
	if err != nil {
		return err
	}
	return err
}

func (p PgxDB) Reserve(userId, serviceId, orderId uint64, amount float32) error {
	ctx := context.TODO()

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

	var user models.User
	err = tx.QueryRow(ctx, "select * from users where id = $1;", userId).Scan(&user.ID, &user.Balance)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("db: reserve: no such user with id %d", userId)
		return err
	} else if err != nil {
		return err
	} else if amount > user.Balance {
		err = fmt.Errorf("db: reserve: the user %d doesn't have enough money, needed: %.2f, user has: %.2f", userId, amount, user.Balance)
		return err
	}

	var checkServiceId uint64
	err = tx.QueryRow(ctx, "select id from services where id = $1;", serviceId).Scan(&checkServiceId)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("db: reserve: no such service with id %d", serviceId)
		return err
	} else if err != nil {
		return err
	}

	user.Balance -= amount
	var updateId uint64
	err = tx.QueryRow(ctx, "update users SET balance = $2 where id = $1 returning id", user.ID, user.Balance).Scan(&updateId)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Europe/Moscow")
	date := time.Now().In(loc)

	var reserveId uint64
	err = tx.QueryRow(ctx, "insert into reserves (order_id, user_id, service_id, amount, purchased, reserved_at) values ($1, $2, $3, $4, $5, $6) returning order_id",
		orderId, userId, serviceId, amount, false, date).Scan(&reserveId)
	if err != nil {
		return err
	}
	return err
}

func (p PgxDB) GetReserve(userId, serviceId, orderId uint64) (models.Reserve, error) {
	ctx := context.TODO()

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
		err = fmt.Errorf("db: get reserve: money were not reserved order %d", orderId)
		return models.Reserve{}, err
	} else if err != nil {
		return models.Reserve{}, err
	}

	var user models.User
	err = tx.QueryRow(ctx, "select * from users where id = $1;", userId).Scan(&user.ID, &user.Balance)
	if err != nil {
		return models.Reserve{}, err
	}
	reserve.User = user

	var service models.Service
	err = tx.QueryRow(ctx, "select * from services where id = $1;", serviceId).Scan(&service.ID, &service.Name)
	if err != nil {
		return models.Reserve{}, err
	}
	reserve.Service = service

	return reserve, err
}

func (p PgxDB) Purchase(userId, serviceId, orderId uint64, amount float32) error {
	ctx := context.TODO()

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

	err = tx.QueryRow(ctx, "select * from reserves where order_id = $1 and user_id = $2 and service_id = $3",
		reserve.OrderID, reserve.UserID, reserve.ServiceID).Scan(&reserve.OrderID, &reserve.UserID, &reserve.ServiceID,
		&reserve.Amount, &reserve.Purchased, &reserve.ReservedAt, &reserve.PurchasedAt)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("db: get purchase: money were not reserved for order %d, user %d and service %d", userId, serviceId, orderId)
		return err
	} else if err != nil {
		return err
	} else if reserve.Amount != amount {
		err = fmt.Errorf("db: purchase: wrong purchase amount, stored in reserve: %.2f, got: %.2f", reserve.Amount, amount)
		return err
	} else if reserve.Purchased {
		err = errors.New("db: purchase: the purchase has already happened")
		return err
	}

	loc, _ := time.LoadLocation("Europe/Moscow")
	purchasedAt := time.Now().In(loc)
	reserve.PurchasedAt = &purchasedAt
	reserve.Purchased = true

	var updateId uint64
	err = tx.QueryRow(ctx, "update reserves SET purchased = $1, purchased_at = $2 where order_id = $3 returning order_id",
		reserve.Purchased, reserve.PurchasedAt, reserve.OrderID).Scan(&updateId)
	if err != nil {
		return err
	}

	return err
}

func (p PgxDB) AddServices(services []models.Service) error {
	var sb strings.Builder
	sb.WriteString("insert into services (id, name) values ")
	for i, s := range services {
		var row string
		if i == 0 {
			row = "(" + strconv.FormatUint(s.ID, 10) + ", '" + s.Name + "')"
		} else {
			row = ", (" + strconv.FormatUint(s.ID, 10) + ", '" + s.Name + "')"
		}
		sb.WriteString(row)
	}

	ctx := context.TODO()

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

	res, err := tx.Exec(ctx, sb.String())
	if err != nil {
		return err
	}
	if res.RowsAffected() != int64(len(services)) {
		err = errors.New("db: add services: failed to add services")
		return err
	}
	return err
}

func (p PgxDB) GetService(serviceId uint64) (models.Service, error) {
	ctx := context.TODO()

	var service models.Service
	err := p.QueryRow(ctx, "select * from services where id = $1;", serviceId).Scan(&service.ID, &service.Name)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("db: reserve: no such service with id %d", serviceId)
		return models.Service{}, err
	} else if err != nil {
		return models.Service{}, err
	}

	return service, err
}
