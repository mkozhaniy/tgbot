package interfaces

import (
	"github.com/tgbot/internal/domain/models"
)

type ProductStorage interface {
	EntityStorage
	GetProductById(id uint) (*models.Product, error)
	GetPruductByName(name string) (*models.Product, error)
	GetProductsByKind(kind string) ([]models.Product, error)
	GetProductByKind(kind string) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
	GetAllKinds() ([]string, error)
}
