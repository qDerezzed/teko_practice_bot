package bot

import (
	"fmt"
	"log"
	"strconv"
	"teko_pracrice_bot/store"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) handlerError(chatID int64, err error) {
	log.Println(err.Error())
	if err = bot.sendMsg(chatID, "Что-то пошло не так. Введите команду /start"); err != nil {
		log.Println(err.Error())
	}
}

func (bot *Bot) handlerCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case "start":
		user := store.NewUser(message.From.ID, store.Start, message.From.UserName)
		if err := bot.storage.AddUser(user); err != nil {
			return err
		}

		if err := bot.printAllCompanys(message.Chat.ID); err != nil {
			return err
		}

		err := bot.sendMsgWithKeyboard(message.Chat.ID, "Добавьте компанию или выберите имеющуюся", StartKeyboard)
		return err
	default:
		err := bot.sendMsg(message.Chat.ID, "Неизвестная команда")
		return err
	}
}

func (bot *Bot) handlerMessage(message *tgbotapi.Message) error {
	chatPos, err := bot.storage.GetChatPosition(message.From.ID)
	if err != nil {
		return err
	}

	switch chatPos {
	case store.Start:
		err := bot.handlerStartPosition(message)
		return err
	case store.EnterCompanyName:
		err := bot.handlerEnterCompanyName(message)
		return err
	case store.EnterCompanyID:
		err := bot.handlerEnterCompanyID(message)
		return err
	case store.CompanyMenu:
		err := bot.handlerCompanyMenu(message)
		return err
	case store.EnterSumExpense:
		err := bot.handlerEnterSumExpense(message)
		return err
	case store.EnterDescExpense:
		err := bot.handlerEnterDescExpense(message)
		return err
	case store.EnterEmailExpense:
		err := bot.handlerEnterEmailExpense(message)
		return err
	}

	return nil
}

func (bot *Bot) handlerStartPosition(message *tgbotapi.Message) error {
	switch message.Text {
	case "Добавить компанию":
		err := bot.storage.SetChatPosition(store.EnterCompanyName, message.From.ID)
		bot.sendMsgWithKeyboard(message.Chat.ID, "Введите название компании", BackKeyboard)
		return err
	case "Выбрать компанию":
		err := bot.storage.SetChatPosition(store.EnterCompanyID, message.From.ID)
		bot.sendMsgWithKeyboard(message.Chat.ID, "Введите ID компании", BackKeyboard)
		return err
	default:
		err := bot.sendMsgWithKeyboard(message.Chat.ID, "Добавьте компанию или выберите имеющуюся", StartKeyboard)
		return err
	}
}

func (bot *Bot) handlerEnterCompanyName(message *tgbotapi.Message) error {
	if err := bot.storage.SetChatPosition(store.Start, message.From.ID); err != nil {
		return err
	}
	if message.Text == "Назад" {
		if err := bot.printAllCompanys(message.Chat.ID); err != nil {
			return err
		}
		err := bot.sendMsgWithKeyboard(message.Chat.ID, "Добавьте компанию или выберите имеющуюся", StartKeyboard)
		return err
	}

	company := store.NewCompany(message.From.ID, message.Text, time.Now())
	if err := bot.storage.AddCompany(company); err != nil {
		return err
	}

	if err := bot.sendMsgWithKeyboard(message.Chat.ID, "Компания успешно добавлена", StartKeyboard); err != nil {
		return err
	}

	err := bot.printAllCompanys(message.Chat.ID)
	return err
}

func (bot *Bot) handlerEnterCompanyID(message *tgbotapi.Message) error {
	if message.Text == "Назад" {
		if err := bot.storage.SetChatPosition(store.Start, message.From.ID); err != nil {
			return err
		}

		if err := bot.printAllCompanys(message.Chat.ID); err != nil {
			return err
		}
		err := bot.sendMsgWithKeyboard(message.Chat.ID, "Добавьте компанию или выберите имеющуюся", StartKeyboard)
		return err
	}

	companyID, convErr := strconv.Atoi(message.Text)
	if convErr != nil {
		err := bot.sendMsgWithKeyboard(message.Chat.ID, "Ошибка: необходимо ввести число", BackKeyboard)
		return err
	}

	isValid, err := bot.storage.IsValidCompanyID(companyID)
	if err != nil || !isValid {
		err := bot.sendMsgWithKeyboard(message.Chat.ID, "Ошибка: компании с таким ID не существует", BackKeyboard)
		return err
	}

	if err := bot.storage.SetChatPosition(store.CompanyMenu, message.From.ID); err != nil {
		return err
	}

	if err := bot.storage.SetCurrentCompanyID(companyID, message.From.ID); err != nil {
		return err
	}

	msg := fmt.Sprintf("Выбрана компания с ID: %d\nНажмите на одну из кнопок", companyID)
	err = bot.sendMsgWithKeyboard(message.Chat.ID, msg, CompanyMenuKeyboard)

	return err
}

func (bot *Bot) handlerCompanyMenu(message *tgbotapi.Message) error {
	switch message.Text {
	case "Создать счёт":
		err := bot.handlerCreateExpense(message)
		return err
	case "Посмотреть последние 10 операций":
		err := bot.handlerShowLastOperations(message)
		return err
	case "Получить информацию по общей сумме платежей за день":
		err := bot.handlerShowInfoPerDay(message)
		return err
	case "В начало":
		if err := bot.storage.SetChatPosition(store.Start, message.From.ID); err != nil {
			return err
		}
		if err := bot.printAllCompanys(message.Chat.ID); err != nil {
			return err
		}
		err := bot.sendMsgWithKeyboard(message.Chat.ID, "Добавьте компанию или выберите имеющуюся", StartKeyboard)
		return err
	default:
		err := bot.sendMsgWithKeyboard(message.Chat.ID, "Нажмите на одну из кнопок", CompanyMenuKeyboard)
		return err
	}
}
