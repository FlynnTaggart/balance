package models

import "time"

type User struct {
	ID      uint64 `json:"id" gorm:"primaryKey"`
	Balance int64  `json:"balance"`
}

type Service struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type Reserve struct {
	OrderID     uint64     `json:"order_id" gorm:"primaryKey;autoIncrement:false"`
	UserID      uint64     `json:"-"`
	User        User       `json:"user"`
	ServiceID   uint64     `json:"-"`
	Service     Service    `json:"service"`
	Amount      int64      `json:"amount"`
	Purchased   bool       `json:"purchased"`
	ReservedAt  time.Time  `json:"reserved_at"`
	PurchasedAt *time.Time `json:"purchased_at,omitempty"` // nullable
}
