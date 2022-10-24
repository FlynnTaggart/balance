package models

import "time"

type User struct {
	ID      uint64  `json:"id" gorm:"primaryKey"`
	Balance float32 `json:"balance"`
}

type Service struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type Order struct {
	ID uint64 `json:"id" gorm:"primaryKey"`
}

type Reserve struct {
	UserID      uint64     `json:"-" gorm:"primaryKey;autoIncrement:false"`
	User        User       `json:"user"`
	ServiceID   uint64     `json:"-" gorm:"primaryKey;autoIncrement:false"`
	Service     Service    `json:"service"`
	OrderID     uint64     `json:"-" gorm:"primaryKey;autoIncrement:false"`
	Order       Order      `json:"order"`
	Amount      float32    `json:"amount"`
	Purchased   bool       `json:"purchased"`
	ReservedAt  time.Time  `json:"reserved_at"`
	PurchasedAt *time.Time `json:"purchased_at,omitempty"` // nullable
}
