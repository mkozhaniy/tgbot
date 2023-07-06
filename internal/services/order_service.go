package services

import (
	"fmt"

	"github.com/tgbot/internal/domain/interfaces"
	"github.com/tgbot/internal/domain/models"
)

func MakeOrder(bucket_storage interfaces.BucketStorage,
	order_storage interfaces.OrderStorage, user_tgid int64) (*models.Order, error) {
	bucket, err := bucket_storage.GetBucketByUserTgid(user_tgid)
	if err != nil {
		return nil, fmt.Errorf("Bucket with tgid %d not found", user_tgid)
	}
	var order models.Order
	order.UserTgid = user_tgid
	order.Products = bucket.Products
	order.Stage = models.Stg(0)
	order.Cost = GetSumCost(*bucket)
	order.Weight = GetSumWeight(*bucket)
	order1, err := order_storage.Save(order)
	if err != nil {
		return nil, fmt.Errorf("Order cannot saved")
	}
	order = order1.(models.Order)
	return &order, nil
}

func DeleteOrder(order_storage interfaces.OrderStorage,
	id uint, user models.User) error {
	if user.Admin {
		order, err := order_storage.GetOrderById(id)
		if err != nil {
			return fmt.Errorf("Order with tgid %d not found", id)
		}
		_, err = order_storage.Delete(order)
		if err != nil {
			return fmt.Errorf("Order with id %d cannot delete", id)
		}
		return nil
	} else {
		orders := user.Orders
		if len(orders) == 0 {
			return fmt.Errorf("User with tgid %d have no orders", user.Tgid)
		}
		for _, order := range orders {
			if order.ID == id {
				_, err := order_storage.Delete(order)
				if err != nil {
					return fmt.Errorf("Order with tgid %d cannot delete", user.Tgid)
				} else {
					return nil
				}
			}
		}
		return fmt.Errorf("Order with id %d in orders of user %d not found",
			id, user.Tgid)
	}
}

func GetStageOrder(order models.Order) string {
	switch order.Stage {
	case models.DURING:
		return "Не оплачен"
	case models.PAYED:
		return "Оплачен"
	case models.DELIVERED:
		return "Доставлен"
	}
	return ""
}
