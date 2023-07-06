package models

type Bucket struct {
	ID       uint      `gorm:"primaryKey"`
	UserTgid int64     `gorm:"unique"`
	Products []Product `gorm:"many2many:bucket_products;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type BucketProduct struct {
	ID        uint `gorm:"primaryKey"`
	BucketID  uint
	ProductID uint
}
