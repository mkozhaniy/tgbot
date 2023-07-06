package interfaces

import (
	"github.com/tgbot/internal/domain/models"
)

type OrderStorage interface {
	EntityStorage
	GetOrderById(id uint) (*models.Order, error)
	GetOrdersByUserTgid(user_tgid int64) ([]models.Order, error)
	GetAllOrders() ([]models.Order, error)
}
