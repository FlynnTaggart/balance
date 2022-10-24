package databases

import (
	"balance/internal/models"

	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

// GetBalance returns balance from user by given id
func (p PgxDB) GetBalance(id uint64) (int64, error) {
	ctx := context.Background()

	var err error
	defer func() {
		if err != nil {
			p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: add balance: %v", err), nil)
		}
	}()

	var balance int64
	err = p.QueryRow(ctx, "select balance from users where id = $1;", id).Scan(&balance)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("db: reserve: no such user with id %d", id)
	} else if err != nil {
		return 0, err
	}
	return balance, err
}

// AddBalance adds money balance of user by given id
// Also writes report to operations table
func (p PgxDB) AddBalance(id uint64, amount int64) error {
	ctx := context.Background()

	var err error
	defer func() {
		if err != nil {
			p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: add balance: %v", err), nil)
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

	var user models.User
	err = tx.QueryRow(ctx, "select * from users where id = $1;", id).Scan(&user.ID, &user.Balance)
	if err != nil && errors.Is(err, pgx.ErrNoRows) { // if the user was not found then create him and add balance
		user = models.User{
			ID:      id,
			Balance: amount,
		}
		var createdId uint64
		err = tx.QueryRow(ctx, "insert into users (id, balance) values ($1, $2) returning id", user.ID, user.Balance).Scan(&createdId)
		if err != nil {
			return err
		}

		loc, _ := time.LoadLocation("Europe/Moscow")
		date := time.Now().In(loc)
		err = tx.QueryRow(ctx, "insert into operations (user_id, amount, done_at) values ($1, $2, $3) returning id", user.ID, user.Balance, date).Scan(&createdId)
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

// DeleteUser deletes user. User cannot be deleted if referenced in reports or operations tables
func (p PgxDB) DeleteUser(id uint64) error {
	ctx := context.Background()

	var err error
	defer func() {
		if err != nil {
			p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: delete user: %v", err), nil)
		}
	}()

	res, err := p.Exec(ctx, "delete from users where id = $1", id)
	if err != nil {
		return err
	} else if res.RowsAffected() == 0 {
		err = fmt.Errorf("db: delete user: no such user with id %d", id)
		return err
	}
	return nil
}
