package tgbot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tgbot/internal/domain/interfaces"
	"github.com/tgbot/internal/domain/models"
	"github.com/tgbot/internal/services"
	"golang.org/x/exp/slices"
)

func addProduct(args string, product_storage interfaces.ProductStorage) (*map[uint]KeyBoardProductInfo, error) {

	rows := len(keyBoardKinds.InlineKeyboard)
	cols := len(keyBoardKinds.InlineKeyboard[0])
	kinds, err := product_storage.GetAllKinds()
	if err != nil {
		return getNavBarProd(product_storage), err
	}
	product, err := services.SaveProduct(product_storage, args)
	if err != nil {
		return getNavBarProd(product_storage), err
	}
	flag := slices.Contains(kinds, product.Kind)
	if !flag {
		if len(keyBoardKinds.InlineKeyboard[rows-1]) <= cols || (rows == 1 && cols <= 3) {
			if len(keyBoardKinds.InlineKeyboard[rows-1]) == 0 {
				keyBoardKinds.InlineKeyboard[rows-1] = make([]tgbotapi.InlineKeyboardButton, 0, 3)
			}
			keyBoardKinds.InlineKeyboard[rows-1] = append(keyBoardKinds.InlineKeyboard[rows-1],
				tgbotapi.NewInlineKeyboardButtonData(product.Kind, "get;"+
					strconv.FormatInt(int64(product.ID), 10)))
		} else {
			keyBoardKinds.InlineKeyboard = append(keyBoardKinds.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(product.Kind,
					"get;"+strconv.FormatInt(int64(product.ID), 10))))
		}
	}
	return getNavBarProd(product_storage), nil
}

func getNavBarProd(product_storage interfaces.ProductStorage) *map[uint]KeyBoardProductInfo {
	products := services.GetProducts(product_storage)
	result := make(map[uint]KeyBoardProductInfo)
	for _, prod := range products {
		n := len(prod)
		if n > 1 {
			result[prod[0].ID] = KeyBoardProductInfo{
				Buttons: getButtonInfo(prod[0],
					prod[n-1].ID, prod[1].ID),
				Prod: prod[0],
			}
			result[prod[n-1].ID] = KeyBoardProductInfo{
				Buttons: getButtonInfo(prod[n-1],
					prod[n-2].ID, prod[0].ID),
				Prod: prod[n-1],
			}
		} else {
			result[prod[0].ID] = KeyBoardProductInfo{
				Buttons: getButtonInfo(prod[0],
					prod[0].ID, prod[0].ID),
				Prod: prod[0],
			}
		}
		for i := 1; i < n-1; i++ {
			result[prod[i].ID] = KeyBoardProductInfo{
				Buttons: getButtonInfo(prod[i], prod[i-1].ID, prod[i+1].ID),
				Prod:    prod[i],
			}
		}
	}
	return &result
}

func fillKindButtons(product_storage interfaces.ProductStorage, kinds []string) {
	cols := 3
	rows := int(len(kinds)/cols + 1)
	buttons := make([][]tgbotapi.InlineKeyboardButton, rows)
	for i := 0; i < rows-1; i++ {
		buttons[i] = make([]tgbotapi.InlineKeyboardButton, cols)
		for j := 0; j < cols; j++ {
			val := kinds[i*cols+j]
			product, err := product_storage.GetProductByKind(val)
			if err != nil {
				continue
			}
			buttons[i][j] = tgbotapi.NewInlineKeyboardButtonData(val, "get;"+
				strconv.FormatInt(int64(product.ID), 10))
		}
	}
	buttons[rows-1] = make([]tgbotapi.InlineKeyboardButton, len(kinds)%cols)
	for i := 0; i < len(kinds)%cols; i++ {
		val := kinds[(rows-1)*cols+i]
		product, err := product_storage.GetProductByKind(val)
		if err != nil {
			continue
		}
		buttons[rows-1][i] = tgbotapi.NewInlineKeyboardButtonData(val, "get;"+
			strconv.FormatInt(int64(product.ID), 10))
	}
	keyBoardKinds.InlineKeyboard = buttons
}

func makeAdmin(msg tgbotapi.MessageConfig, arg string,
	user_storage interfaces.UserStorage) tgbotapi.MessageConfig {
	args := strings.Split(arg, " ")
	if len(args) > 1 && args[0] == ADMINS_KEY {
		user_tgid, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			msg.Text = "Неверный id пользователя, либо пользователь еще не пользовался ботом"
		}
		if err := services.MakeUserAdmin(user_storage, user_tgid); err != nil {
			msg.Text = "Произошла ошибка на стороне бота("
		}
	} else {
		msg.Text = "Некорректно введены слова после комманды или неверный ключ"
	}
	return msg
}

func addPhoto(bot *tgbotapi.BotAPI, update tgbotapi.Update,
	product_storage interfaces.ProductStorage) error {
	text := update.Message.ReplyToMessage.Text
	err := services.AddPhotosByName(product_storage, strings.Split(text, " ")[1],
		update.Message.Photo[0].FileID)
	return err
}

func getButtonInfo(product models.Product, prev uint, next uint) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину",
				"add_bucket;"+strconv.FormatInt(int64(product.ID), 10)),
			tgbotapi.NewInlineKeyboardButtonData("Фотографии",
				"get_photos;"+strconv.FormatInt(int64(product.ID), 10)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущий", "get;"+
				strconv.FormatInt(int64(prev), 10)),
			tgbotapi.NewInlineKeyboardButtonData("Следующий", "get;"+
				strconv.FormatInt(int64(next), 10)),
		),
	)
}

func viewProduct(msg tgbotapi.MessageConfig,
	navigationBar KeyBoardProductInfo) tgbotapi.MessageConfig {
	msg.Text += fmt.Sprintf("Название - %s\n", navigationBar.Prod.Name)
	msg.Text += fmt.Sprintf("Цена - %.2f руб., \nВес - %.2f кг., \nОписание - %s",
		navigationBar.Prod.Cost, navigationBar.Prod.Weight, navigationBar.Prod.Description)
	msg.ReplyMarkup = navigationBar.Buttons
	return msg
}

func viewBucket(msg tgbotapi.MessageConfig, bucket models.Bucket) tgbotapi.MessageConfig {
	msg.Text = ""
	for _, prod := range bucket.Products {
		msg.Text += fmt.Sprintf("%s, %.2f кг., %.2f руб.\n",
			prod.Name, prod.Weight, prod.Cost)
	}
	msg.Text += fmt.Sprintf("Общая стоимость: %.2f руб., Общий вес: %.2f кг.",
		services.GetSumCost(bucket), services.GetSumWeight(bucket))
	msg.ReplyMarkup = keyBoardBucket
	return msg
}

func viewOrders(msg tgbotapi.MessageConfig, orders []models.Order) tgbotapi.MessageConfig {
	msg.Text = ""
	for _, order := range orders {
		msg.Text += fmt.Sprintf("Номер - %d, Стадия - %s, Стоимость - %.2f руб., Вес - %.2f кг.\n",
			order.ID, services.GetStageOrder(order), order.Cost, order.Weight)
	}
	msg.ReplyMarkup = keyBoardOrders
	return msg
}

func viewOrder(msg tgbotapi.MessageConfig, order_storage interfaces.OrderStorage,
	id uint) tgbotapi.MessageConfig {
	order, err := order_storage.GetOrderById(id)
	if err != nil {
		msg.Text = "Неверный номер заказа"
	} else {
		msg.Text = ""
		for _, prod := range order.Products {
			msg.Text += fmt.Sprintf("%s, Цена - %.2f руб., Вес - %.2f кг.\n",
				prod.Name, prod.Cost, prod.Weight)
		}
		msg.Text += fmt.Sprintf("Общая стоимость: %.2f руб., Общий вес: %.2f кг.",
			order.Cost, order.Weight)
	}
	return msg
}

func viewPhotoProduct(product_storaga interfaces.ProductStorage,
	id uint, chat_id int64) (string, *tgbotapi.MediaGroupConfig) {
	product, err := product_storaga.GetProductById(id)
	if err != nil || len(product.Photo) == 0 {
		return "У данного товара еще нет фотографий", nil
	}

	photos := make([]interface{}, len(product.Photo))
	for i, photo := range product.Photo {
		photos[i] = tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(photo))
	}
	mgrp := tgbotapi.NewMediaGroup(chat_id, photos)
	return "", &mgrp
}

func deleteProduct(product_storage interfaces.ProductStorage,
	name string) (string, *map[uint]KeyBoardProductInfo, error) {
	_, err := services.DeleteProduct(product_storage, name)
	if err != nil {
		return fmt.Sprintf("Товар с название %s не удалось удалить", name),
			getNavBarProd(product_storage), err
	} else {
		kinds, err := product_storage.GetAllKinds()
		if err != nil {
			return "Произошла ошибка на сервере", getNavBarProd(product_storage),
				fmt.Errorf("Error with getting kinds of products")
		}
		fillKindButtons(product_storage, kinds)
		return fmt.Sprintf("Товар с название %s удален", name), getNavBarProd(product_storage), nil
	}

}
