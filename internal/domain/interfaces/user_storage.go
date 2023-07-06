package interfaces

import (
	"github.com/tgbot/internal/domain/models"
)

type UserStorage interface {
	EntityStorage
	GetUserByTgid(user_tgid int64) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
	GetUserById(id uint) (*models.User, error)
}
