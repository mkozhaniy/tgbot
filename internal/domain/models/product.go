package models

import "github.com/lib/pq"

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique"`
	Weight      float32
	Cost        float32
	Amount      uint
	Photo       pq.StringArray `gorm:"type:text[]"`
	Kind        string
	Description string
}
