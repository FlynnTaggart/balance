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
		IsoLevel:       pgx.Serializable,
		DeferrableMode: pgx.Deferrable,
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
	return errors.New("db: method not implemented")
}

func (p PgxDB) Purchase(userId, serviceId, orderId uint64, amount float32) error {
	return errors.New("db: method not implemented")
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
		IsoLevel:       pgx.Serializable,
		DeferrableMode: pgx.Deferrable,
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
		err = errors.New("db: add services: ")
		return err
	}
	return err
}

func (p PgxDB) GetReserve(userId, serviceId, orderId uint64) (models.Reserve, error) {
	return models.Reserve{}, errors.New("db: method not implemented")
}

func (p PgxDB) GetService(serviceId uint64) (models.Service, error) {
	ctx := context.TODO()

	var service models.Service
	err := p.QueryRow(ctx, "select * from services where id = $1;", serviceId).Scan(&service.ID, &service.Name)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return models.Service{}, fmt.Errorf("db: reserve: no such service with id %d", serviceId)
	} else if err != nil {
		return models.Service{}, err
	}

	return service, err
}
