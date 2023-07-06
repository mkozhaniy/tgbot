package interfaces

import (
	"github.com/tgbot/internal/domain/models"
)

type BucketStorage interface {
	EntityStorage
	GetBucketByUserTgid(user_tgid int64) (*models.Bucket, error)
	DeleteProductFromBucket(bucket *models.Bucket, product *models.Product) (*models.Bucket, error)
	ClearBucket(bucket *models.Bucket) (*models.Bucket, error)
	AddProductToBucket(bucket *models.Bucket, product *models.Product) (*models.Bucket, error)
}
