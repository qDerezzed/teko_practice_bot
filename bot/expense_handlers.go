package bot

import (
	"fmt"
	"net/mail"
	"strconv"
	"teko_pracrice_bot/store"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) handlerCreateExpense(message *tgbotapi.Message) error {
	if err := bot.storage.SetChatPosition(store.EnterSumExpense, message.From.ID); err != nil {
		return err
	}
	err := bot.sendMsgWithKeyboard(message.Chat.ID, "Введите сумму", BackKeyboard)
	return err
}

func (bot *Bot) handlerShowLastOperations(message *tgbotapi.Message) error {
	companyID, err := bot.storage.GetCurrentCompanyID(message.From.ID)
	if err != nil {
		return err
	}
	expenses, err := bot.storage.GetLastExpenses(companyID, 10)
	if err != nil {
		return err
	}

	if expenses == nil {
		err := bot.sendMsg(message.Chat.ID, "Операции не совершались")
		return err
	}
	if err := bot.sendMsg(message.Chat.ID, "Последние 10 операций: "); err != nil {
		return err
	}
	for _, expense := range expenses {
		if err := bot.printExpense(expense, message.Chat.ID); err != nil {
			return err
		}
	}

	return err
}

func (bot *Bot) handlerShowInfoPerDay(message *tgbotapi.Message) error {
	companyID, err := bot.storage.GetCurrentCompanyID(message.From.ID)
	if err != nil {
		return err
	}
	sum, err := bot.storage.GetSumPerDay(companyID)
	if err != nil {
		return err
	}

	infoMsg := fmt.Sprintf("Общая сумма платежей в компании за день: %d", sum)

	if err := bot.sendMsg(message.Chat.ID, infoMsg); err != nil {
		return err
	}
	return err
}

func (bot *Bot) handlerEnterSumExpense(message *tgbotapi.Message) error {
	companyID, err := bot.storage.GetCurrentCompanyID(message.From.ID)
	if err != nil {
		return err
	}
	if message.Text == "Назад" {
		if err := bot.storage.SetChatPosition(store.CompanyMenu, message.From.ID); err != nil {
			return err
		}
		msg := fmt.Sprintf("Выбрана компания с ID: %d\nНажмите на одну из кнопок", companyID)
		err := bot.sendMsgWithKeyboard(message.Chat.ID, msg, CompanyMenuKeyboard)
		return err
	}

	sum, convErr := strconv.Atoi(message.Text)
	if convErr != nil {
		err := bot.sendMsgWithKeyboard(message.Chat.ID, "Ошибка: необходимо ввести число", BackKeyboard)
		return err
	}

	expense := store.NewExpense(sum, "", "", time.Now(), companyID, message.From.ID)
	expenseID, err := bot.storage.AddExpense(expense)
	if err != nil {
		return err
	}

	bot.storage.SetCurrentExpenseID(expenseID, message.From.ID)

	if err := bot.storage.SetChatPosition(store.EnterDescExpense, message.From.ID); err != nil {
		return err
	}

	err = bot.sendMsgWithKeyboard(message.Chat.ID, "Введите описание", BackKeyboard)

	return err
}

func (bot *Bot) handlerBackFromExpense(message *tgbotapi.Message, expenseID int) error {
	if err := bot.storage.DeleteExpense(expenseID); err != nil {
		return err
	}
	companyID, err := bot.storage.GetCurrentCompanyID(message.From.ID)
	if err != nil {
		return err
	}
	if err := bot.storage.SetChatPosition(store.CompanyMenu, message.From.ID); err != nil {
		return err
	}
	msg := fmt.Sprintf("Выбрана компания с ID: %d\nНажмите на одну из кнопок", companyID)
	err = bot.sendMsgWithKeyboard(message.Chat.ID, msg, CompanyMenuKeyboard)
	return err
}

func (bot *Bot) handlerEnterDescExpense(message *tgbotapi.Message) error {
	expenseID, err := bot.storage.GetCurrentExpenseID(message.From.ID)
	if err != nil {
		return err
	}

	if message.Text == "Назад" {
		err := bot.handlerBackFromExpense(message, expenseID)
		return err
	}

	if err := bot.storage.SetDescriptionExpense(message.Text, expenseID); err != nil {
		return err
	}

	if err := bot.storage.SetChatPosition(store.EnterEmailExpense, message.From.ID); err != nil {
		return err
	}
	err = bot.sendMsgWithKeyboard(message.Chat.ID, "Введите email", BackKeyboard)

	return err
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (bot *Bot) handlerEnterEmailExpense(message *tgbotapi.Message) error {
	expenseID, err := bot.storage.GetCurrentExpenseID(message.From.ID)
	if err != nil {
		return err
	}

	if message.Text == "Назад" {
		err := bot.handlerBackFromExpense(message, expenseID)
		return err
	}

	if !isEmail(message.Text) {
		err := bot.sendMsgWithKeyboard(message.Chat.ID, "Ошибка: необходимо ввести email", BackKeyboard)
		return err
	}

	if err := bot.storage.SetEmailExpense(message.Text, expenseID); err != nil {
		return err
	}

	if err := bot.storage.SetChatPosition(store.CompanyMenu, message.From.ID); err != nil {
		return err
	}

	err = bot.sendMsgWithKeyboard(message.Chat.ID, "Нажмите на одну из кнопок", CompanyMenuKeyboard)

	return err
}
