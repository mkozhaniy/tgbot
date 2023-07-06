package services

import (
	"fmt"

	"github.com/tgbot/internal/domain/interfaces"
	"github.com/tgbot/internal/domain/models"
)

func AddProductToBucket(bucket_storage interfaces.BucketStorage,
	product_storage interfaces.ProductStorage, user_tgid int64,
	product_id uint) (*models.Bucket, error) {
	bucket, err := bucket_storage.GetBucketByUserTgid(user_tgid)
	if err != nil {
		return nil, fmt.Errorf("Bucket for user with tgid %d not found", user_tgid)
	}
	product, err := product_storage.GetProductById(product_id)
	if err != nil {
		return nil, fmt.Errorf("Product with id %d not found", product_id)
	}
	bucket, err = bucket_storage.AddProductToBucket(bucket, product)
	if err != nil {
		return nil, fmt.Errorf("Cannot add product(id:%d) to bucket(tgid:%d)",
			product_id, user_tgid)
	}
	return bucket, nil
}

func DeleteProductFromBucket(bucket_storage interfaces.BucketStorage,
	product_storage interfaces.ProductStorage,
	user_tgid int64, name string) (*models.Bucket, error) {
	bucket, err := bucket_storage.GetBucketByUserTgid(user_tgid)
	if err != nil {
		return nil, fmt.Errorf("Cannot delete product(name:%s) from bucket(tgid:%d)",
			name, user_tgid)
	}
	product, err := product_storage.GetPruductByName(name)
	if err != nil {
		return nil, fmt.Errorf("Product with name %s not exist", name)
	}
	bucket, err = bucket_storage.DeleteProductFromBucket(bucket, product)
	if err != nil {
		return nil, fmt.Errorf("Cannot delete product(name:%s) from bucket(tgid:%d)",
			name, user_tgid)
	}
	return bucket, nil
}

func GetAllProducts(source interfaces.BucketStorage, user_tgid int64) []models.Product {
	bucket, err := source.GetBucketByUserTgid(user_tgid)
	if err != nil {
		return nil
	}
	return bucket.Products
}

func ClearBucket(source interfaces.BucketStorage, user_tgid int64) (*models.Bucket, error) {
	bucket, err := source.GetBucketByUserTgid(user_tgid)
	if err != nil {
		return nil, fmt.Errorf("Bucket with tgid %d not found", user_tgid)
	}
	bucket, err = source.ClearBucket(bucket)
	if err != nil {
		return nil, fmt.Errorf("Bucket with tgid %d not cleared", user_tgid)
	}
	return bucket, nil
}

func GetSumCost(bucket models.Bucket) float32 {
	var result float32
	for _, val := range bucket.Products {
		result += val.Cost
	}
	return result
}

func GetSumWeight(bucket models.Bucket) float32 {
	var result float32
	for _, val := range bucket.Products {
		result += val.Weight
	}
	return result
}
