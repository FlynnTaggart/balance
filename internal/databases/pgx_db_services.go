package databases

import (
	"balance/internal/models"

	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

// AddServices adds an array of services
func (p PgxDB) AddServices(services []models.Service) error {
	// create one big query with adding all given services
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

	ctx := context.Background()

	var err error
	defer func() {
		if err != nil {
			p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: add services: %v", err), nil)
		}
	}()

	// exec the query
	res, err := p.Exec(ctx, sb.String())
	if err != nil {
		return err
	}
	if res.RowsAffected() != int64(len(services)) {
		err = errors.New("db: add services: failed to add services")
		return err
	}
	return err
}

// GetService returns service by given id
func (p PgxDB) GetService(id uint64) (models.Service, error) {
	ctx := context.Background()

	var err error
	defer func() {
		if err != nil {
			p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: get service: %v", err), nil)
		}
	}()

	var service models.Service
	err = p.QueryRow(ctx, "select * from services where id = $1;", id).Scan(&service.ID, &service.Name)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("db: get service: no such service with id %d", id)
		return models.Service{}, err
	} else if err != nil {
		return models.Service{}, err
	}

	return service, err
}

// DeleteService deletes service by given id. Service cannot be deleted if referenced in reports or operations tables
func (p PgxDB) DeleteService(id uint64) error {
	ctx := context.Background()

	var err error
	defer func() {
		if err != nil {
			p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: delete service: %v", err), nil)
		}
	}()

	res, err := p.Exec(ctx, "delete from services where id = $1", id)
	if err != nil {
		return err
	} else if res.RowsAffected() == 0 {
		err = fmt.Errorf("db: delete service: no such service with id %d", id)
		return err
	}
	return nil
}
