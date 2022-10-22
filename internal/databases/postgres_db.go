package databases

import (
	"balance/internal/models"

	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresDB struct {
	*gorm.DB
}

func NewPostgresDB(DB *gorm.DB) *PostgresDB {
	return &PostgresDB{DB}
}

func (db *PostgresDB) GetBalance(id uint64) (float32, error) {
	var user models.User
	if err := db.Take(&user, id).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("db: reserve: no such user with id %d", id)
	} else if err != nil {
		return 0, err
	}
	return user.Balance, nil
}

func (db *PostgresDB) AddBalance(id uint64, amount float32) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	err := tx.Take(&user, id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		user = models.User{
			ID:      id,
			Balance: amount,
		}
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return err
		}

		return tx.Commit().Error

	} else if err != nil {
		tx.Rollback()
		return err
	}

	user.Balance += amount
	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (db *PostgresDB) Reserve(userId, serviceId, orderId uint64, amount float32) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	err := tx.Take(&user, userId).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return fmt.Errorf("db: reserve: no such user with id %d", userId)
	} else if err != nil {
		tx.Rollback()
		return err
	}

	if amount > user.Balance {
		tx.Rollback()
		return fmt.Errorf("db: reserve: the user %d doesn't have enough money, needed: %.2f, user has: %.2f", userId, amount, user.Balance)
	}

	var service models.Service
	err = tx.Take(&service, serviceId).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return fmt.Errorf("db: reserve: no such service with id %d", serviceId)
	} else if err != nil {
		tx.Rollback()
		return err
	}

	user.Balance -= amount
	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	order := models.Order{ID: orderId}

	reserve := models.Reserve{
		User:       user,
		Service:    service,
		Order:      order,
		Amount:     amount,
		ReservedAt: time.Now(),
	}
	tx.Clauses(clause.OnConflict{DoNothing: true})
	err = tx.Create(&reserve).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (db *PostgresDB) Purchase(userId, serviceId, orderId uint64, amount float32) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	reserve := models.Reserve{
		UserID:    userId,
		ServiceID: serviceId,
		OrderID:   orderId,
	}
	err := tx.Preload("User").Preload("Service").Preload("Order").Take(&reserve).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return fmt.Errorf("db: purchase: money were not reserved for user %d, service %d and order %d", userId, serviceId, orderId)
	} else if err != nil {
		tx.Rollback()
		return err
	} else if reserve.Amount != amount {
		tx.Rollback()
		return fmt.Errorf("db: purchase: wrong purchase amount, stored in reserve: %.2f, got: %.2f", reserve.Amount, amount)
	}

	report := models.Report{
		Service:     reserve.Service,
		Amount:      reserve.Amount,
		PurchasedAt: time.Now(),
	}

	tx.Clauses(clause.OnConflict{DoNothing: true})
	err = tx.Create(&report).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Delete(&reserve).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (db *PostgresDB) AddServices(services []models.Service) error {
	return db.Create(&services).Error
}

func (db *PostgresDB) GetReserve(userId, serviceId, orderId uint64) (models.Reserve, error) {
	reserve := models.Reserve{
		UserID:    userId,
		ServiceID: serviceId,
		OrderID:   orderId,
	}
	err := db.Preload("User").Preload("Service").Preload("Order").Take(&reserve).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Reserve{}, fmt.Errorf("db: get reserve: money were not reserved for user %d, service %d and order %d", userId, serviceId, orderId)
	} else if err != nil {
		return models.Reserve{}, err
	}
	return reserve, nil
}

func (db *PostgresDB) GetService(serviceId uint64) (models.Service, error) {
	var service models.Service
	err := db.Take(&service, serviceId).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Service{}, fmt.Errorf("db: reserve: no such service with id %d", serviceId)
	} else if err != nil {
		return models.Service{}, err
	}
	return service, nil
}
