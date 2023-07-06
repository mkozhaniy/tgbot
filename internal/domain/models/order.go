package models

import (
	"time"
)

type Stg int

const (
	DURING Stg = iota
	PAYED
	DELIVERED
)

type Order struct {
	ID        uint `gorm:"primaryKey"`
	UserTgid  int64
	Stage     Stg
	Products  []Product `gorm:"many2many:order_products;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Cost      float32
	Weight    float32
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Payed     time.Time
	Delivered time.Time
}

type OrderProduct struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	ProductID uint
}
