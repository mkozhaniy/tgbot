package tgbot

import (
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tgbot/internal/domain/interfaces"
	"github.com/tgbot/internal/services"
	"go.uber.org/zap"
)

var ADMINS_KEY = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
var logger1 *zap.Logger

func stop_bot(bot *tgbotapi.BotAPI, ch chan struct{}, logger *zap.Logger) {
	<-ch
	bot.StopReceivingUpdates()
	logger.Sugar().Info("bot is inerupted")
	logger.Sync()
}

func Start(logger *zap.Logger, product_storage interfaces.ProductStorage,
	user_storage interfaces.UserStorage, bucket_storage interfaces.BucketStorage,
	order_storage interfaces.OrderStorage, token string, ch chan struct{}) error {
	gap := 60
	bot, err := tgbotapi.NewBotAPI(token)
	logger1 = logger
	go stop_bot(bot, ch, logger)
	if err != nil {
		panic(err.Error())
	}

	timer := time.NewTimer(time.Duration(gap) * time.Second)
	go func(logger *zap.Logger) {
		<-timer.C
		logger.Sync()
	}(logger)

	logger.Sugar().Infof("Authorized on account %s", bot.Self.UserName)

	kinds, err := product_storage.GetAllKinds()
	if err != nil {
		panic(err)
	}
	fillKindButtons(product_storage, kinds)
	navBarProduct := getNavBarProd(product_storage)

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Message.Text)
			if _, err := bot.Request(callback); err != nil {
				logger.DPanic(err.Error())
			}

			delMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID)
			_, err := bot.Send(delMsg)
			if err != nil {
				logger.Sugar().Info(err.Error())
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
			switch str := strings.Split(update.CallbackQuery.Data, ";"); str[0] {
			case "get":
				id, err := strconv.ParseInt(str[1], 10, 64)
				if err != nil {
					logger.Sugar().Info(err.Error())
					msg.Text = "Произошла ошибка на сервере"
					continue
				}
				msg = viewProduct(msg, (*navBarProduct)[uint(id)])
			case "add_bucket":
				id, err := strconv.ParseInt(str[1], 10, 64)
				if err != nil {
					logger.Sugar().Info(err.Error())
					msg.Text = "Произошла ошибка на сервере"
					continue
				}

				_, err = services.AddProductToBucket(bucket_storage, product_storage,
					update.CallbackQuery.From.ID, uint(id))
				if err != nil {
					msg.Text = "Товар не удалось добавить в корзину"
					logger.Sugar().Info(err.Error())
				} else {
					msg.Text = "Товар добвален в корзину"
				}
			case "get_photos":
				id, err := strconv.ParseInt(str[1], 10, 64)
				if err != nil {
					logger.Sugar().Info(err.Error())
					msg.Text = "Произошла ошибка на сервере"
					continue
				}
				text, mgrp := viewPhotoProduct(product_storage, uint(id), msg.ChatID)
				msg.Text = text
				if mgrp != nil {
					bot.SendMediaGroup(*mgrp)
				}
			case "make_order":
				_, err := services.MakeOrder(bucket_storage, order_storage, update.CallbackQuery.From.ID)
				if err == nil {
					msg.Text = "Заказ оформлен"
				} else {
					logger.Sugar().Info(err.Error())
					msg.Text = "Произошла ошибка на стороне сервера"
				}
			case "clear_bucket":
				_, err := services.ClearBucket(bucket_storage, update.CallbackQuery.From.ID)
				if err != nil {
					msg.Text = "Произошла ошибка на стороне сервера"
					logger.Sugar().Info(err.Error())
				} else {
					msg.Text = "Корзина очищена"
				}
			case "delete_product":
				msg.Text = "Напишите в этот чат \"delete_product название_товара\""
			case "info_order":
				msg.Text = "Напишите в этот чат \"info_order номер_заказа\""
			case "delete_order":
				msg.Text = "Напишите в этот чат \"delete_order номер_заказа\""
			}
			_, err = bot.Send(msg)
			if err != nil {
				logger.Sugar().Info(err.Error())
			}

		} else if update.Message != nil { // If we got a messag
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите категорию")
			msg.ReplyMarkup = keyBoardClient
			if usr, err := user_storage.GetUserByTgid(update.Message.From.ID); err == nil && usr.Admin {
				if update.Message.IsCommand() {
					switch cmd := update.Message.Command(); cmd {
					case "add_product":
						navBarProduct, err = addProduct(update.Message.CommandArguments(), product_storage)
						if err != nil {
							logger.Info(err.Error())
							msg.Text = "Товар не удалось добавить"
						} else {
							msg.Text = "Товар успешно добавлен"
						}
					case "make_admin":
						msg = makeAdmin(msg, update.Message.CommandArguments(), user_storage)
					case "delete_prod":
						args := strings.Split(update.Message.CommandArguments(), " ")
						if len(args) == 0 {
							msg.Text = "Вы не написали название товара"
						} else {
							msg.Text, navBarProduct, err = deleteProduct(product_storage, args[0])
							if err != nil {
								logger.Info(err.Error())
							}
						}
					default:
						msg.Text = "Неизвестная комманда"
					}
				} else {
					if update.Message.Photo != nil && update.Message.ReplyToMessage != nil &&
						strings.Split(update.Message.ReplyToMessage.Text, " ")[0] == "/add_product" {
						err := addPhoto(bot, update, product_storage)
						if err != nil {
							logger.Info(err.Error())
							msg.Text = "Не удалось добавить фотографию, неверные данные"
						} else {
							msg.Text = "Фотография успешно добавлена к товару"
						}
					}
				}
			} else if update.Message.IsCommand() {
				switch cmd := update.Message.Command(); cmd {
				case "start":
					services.SaveUser(user_storage, update.SentFrom().ID, update.SentFrom().UserName)
					logger.Sugar().Infof("User %s is registred", update.SentFrom().UserName)
				case "make_admin":
					msg = makeAdmin(msg, update.Message.CommandArguments(), user_storage)
				default:
					msg.Text = "Неизвестная комманда"
				}
			}
			str := strings.Split(update.Message.Text, " ")
			switch strings.Trim(str[0], " ") {
			case "Товары":
				msg.Text = "Выберите категорию товара"
				msg.ReplyMarkup = keyBoardKinds
			case "Корзина":
				bucket, err := bucket_storage.GetBucketByUserTgid(update.Message.From.ID)
				if err == nil {
					msg = viewBucket(msg, *bucket)
				} else {
					logger.Sugar().Info(err.Error())
					msg.Text = "Произошла ошибка на сервере(вашей корзины не найдено)"
				}
			case "Заказы":
				usr, err := user_storage.GetUserByTgid(update.Message.From.ID)
				if usr.Admin {
					orders, err := order_storage.GetAllOrders()
					if err == nil && len(orders) > 0 {
						msg = viewOrders(msg, orders)
					} else {
						msg.Text = "Нет заказов"
					}
				} else {
					if err == nil && len(usr.Orders) > 0 {
						msg = viewOrders(msg, usr.Orders)
					} else {
						msg.Text = "У вас пока нет заказов"
					}
				}
			case "delete_product":
				if len(str) <= 1 {
					msg.Text = "Вы забыли написать название товара"
				} else {
					_, err := services.DeleteProductFromBucket(bucket_storage, product_storage,
						update.Message.From.ID, strings.Trim(str[1], " "))
					if err != nil {
						logger.Sugar().Info(err.Error())
						msg.Text = "Товар не удалось убрать из корзины"
					} else {
						msg.Text = "Товар убран из корзины"
					}
				}
			case "delete_order":
				if len(str) <= 1 {
					msg.Text = "Вы забыли написать номер заказа"
				} else {
					id, err := strconv.ParseInt(strings.Trim(str[1], " "), 10, 64)
					if err != nil {
						msg.Text = "Неверный номер заказа"
					} else {
						usr, err := user_storage.GetUserByTgid(update.Message.From.ID)
						if err != nil {
							logger.Sugar().
								Info("User with tgid %d not found", update.Message.From.ID)
						}
						err = services.DeleteOrder(order_storage,
							uint(id), *usr)
						if err != nil {
							logger.Sugar().Info(err.Error())
							msg.Text = "Заказ не удалось удалить"
						} else {
							msg.Text = "Заказ успешно удален"
						}
					}
				}
			case "info_order":
				if len(str) <= 1 {
					msg.Text = "Вы забыли написать номер заказа"
				} else {
					id, err := strconv.ParseInt(strings.Trim(str[1], " "), 10, 64)
					if err != nil {
						logger.Sugar().Info(err.Error())
						msg.Text = "Неверный номер заказа"
					} else {
						msg = viewOrder(msg, order_storage, uint(id))
					}
				}
			}
			_, err = bot.Send(msg)
			if err != nil {
				logger.Sugar().Info(err.Error())
			}
		}
	}
	return nil
}
