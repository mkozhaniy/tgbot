package services

import (
	"fmt"

	"github.com/tgbot/internal/domain/interfaces"
	"github.com/tgbot/internal/domain/models"
)

func SaveUser(source interfaces.UserStorage, user_tgid int64, name string) {
	user := models.User{
		Tgid:     user_tgid,
		Username: name,
		Bucket: models.Bucket{
			UserTgid: user_tgid,
		},
		Orders: []models.Order{},
	}
	source.Save(user)
}

func DeleteUser(source interfaces.UserStorage, name string) error {
	user, err := source.GetUserByName(name)
	if err != nil {
		return fmt.Errorf("User with name %s not deleted", name)
	}
	source.Delete(user)
	return nil
}

func MakeUserAdmin(source interfaces.UserStorage, user_tgid int64) error {
	usr, err := source.GetUserByTgid(user_tgid)
	if err != nil {
		return fmt.Errorf("User with tgid %d not found", user_tgid)
	}
	usr.Admin = true
	source.Save(*usr)
	return nil
}
