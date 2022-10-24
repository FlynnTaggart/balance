package databases

import "balance/internal/models"

type DBInt interface {
	GetBalance(id uint64) (float32, error)
	AddBalance(id uint64, amount float32) error
	DeleteUser(id uint64) error
	Reserve(userId, serviceId, orderId uint64, amount float32) error
	GetReserve(userId, serviceId, orderId uint64) (models.Reserve, error)
	DeleteReserve(userId, serviceId, orderId uint64, amount float32) error
	Purchase(userId, serviceId, orderId uint64, amount float32) error
	AddServices(services []models.Service) error
	GetService(id uint64) (models.Service, error)
	DeleteService(id uint64) error
}
