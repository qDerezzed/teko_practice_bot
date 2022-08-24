package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var BackKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Назад"),
	),
)

var StartKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить компанию"),
		tgbotapi.NewKeyboardButton("Выбрать компанию"),
	),
)

var CompanyMenuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Создать счёт"),
		tgbotapi.NewKeyboardButton("Посмотреть последние 10 операций"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Получить информацию по общей сумме платежей за день"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("В начало"),
	),
)
