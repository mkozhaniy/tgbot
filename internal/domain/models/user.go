package models

type User struct {
	ID       uint    `gorm:"primaryKey"`
	Tgid     int64   `gorm:"unique"`
	Username string  `gorm:"unique"`
	Bucket   Bucket  `gorm:"foreignKey:UserTgid;references:Tgid;constraint:OnUpdate:CASCADE;"`
	Orders   []Order `gorm:"foreignKey:UserTgid;references:Tgid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Admin    bool
}
