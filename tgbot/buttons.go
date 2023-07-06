package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tgbot/internal/domain/models"
)

var keyBoardClient = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Товары"),
		tgbotapi.NewKeyboardButton("Заказы"),
		tgbotapi.NewKeyboardButton("Корзина"),
	),
)

var keyBoardKinds = tgbotapi.NewInlineKeyboardMarkup()

var keyBoardBucket = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Оформить заказ", "make_order"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить товар", "delete_product"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Очистить корзину", "clear_bucket"),
	),
)

var keyBoardOrders = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Подробней о заказе", "info_order"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Оплатить заказ", "pay_order"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить заказ", "delete_order"),
	),
)

type KeyBoardProductInfo struct {
	Prod    models.Product
	Buttons tgbotapi.InlineKeyboardMarkup
}
