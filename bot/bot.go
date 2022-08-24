package bot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"teko_pracrice_bot/store"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	storage *store.Store
	botAPI  *tgbotapi.BotAPI
}

func New(storage *store.Store, bot *tgbotapi.BotAPI) *Bot {
	return &Bot{
		storage: storage,
		botAPI:  bot,
	}
}

func GetAPI(botToken string) (*tgbotapi.BotAPI, error) {
	botAPI, err := tgbotapi.NewBotAPI(botToken)
	return botAPI, err
}

func (bot *Bot) SetWebHook(webHookURL string) error {
	webHook, err := tgbotapi.NewWebhook(webHookURL)
	if err != nil {
		return err
	}

	_, err = bot.botAPI.Request(webHook)

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all is working"))
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()

	return err
}

func (bot *Bot) sendMsg(chatID int64, textMsg string) error {
	msg := tgbotapi.NewMessage(chatID, textMsg)
	_, err := bot.botAPI.Send(msg)
	return err
}

func (bot *Bot) sendMsgWithKeyboard(chatID int64, textMsg string, keyboard tgbotapi.ReplyKeyboardMarkup) error {
	startMsg := tgbotapi.NewMessage(chatID, textMsg)
	startMsg.ReplyMarkup = keyboard
	_, err := bot.botAPI.Send(startMsg)
	return err
}

func (bot *Bot) processUpdates(updatesChannel *tgbotapi.UpdatesChannel) error {
	for update := range *updatesChannel {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := bot.handlerCommand(update.Message); err != nil {
				bot.handlerError(update.Message.Chat.ID, err)
			}
			continue
		}

		if err := bot.handlerMessage(update.Message); err != nil {
			bot.handlerError(update.Message.Chat.ID, err)
		}
	}
	return nil
}

func (bot *Bot) Start(ctx context.Context) error {
	// bot.botAPI.Debug = true

	log.Printf("Authorized on account %s", bot.botAPI.Self.UserName)

	updates := bot.botAPI.ListenForWebhook("/")

	err := bot.processUpdates(&updates)

	return err
}

func (bot *Bot) printCompany(company *store.Company, chatID int64) error {
	// userName, err := bot.storage.GetUserName(company.UserID)
	// if err != nil {
	// 	return err
	// }

	msg := fmt.Sprintf(
		"ID: %d\nНазвание: %s\nДата добавления: %s\nДобавлена пользователем с ID: %d",
		company.CompanyID,
		company.Name,
		company.CreationDate.Format("2006-01-02"),
		// userName,
		company.UserID)
	err := bot.sendMsg(chatID, msg)
	return err
}

func (bot *Bot) printAllCompanys(chatID int64) error {
	companys, err := bot.storage.GetAllCompanys()
	if err != nil {
		return err
	}

	if err := bot.sendMsg(chatID, "Список добавленных компаний:"); err != nil {
		return err
	}
	for _, company := range companys {
		if err := bot.printCompany(company, chatID); err != nil {
			return err
		}
	}

	return nil
}

func (bot *Bot) printExpense(expense *store.Expense, chatID int64) error {
	// userName, err := bot.storage.GetUserName(expense.UserID)
	// if err != nil {
	// 	return err
	// }

	expenseMsg := fmt.Sprintf(
		"Сумма: %d\nОписание: %s\nEmail: %s\nДата и время добавления: %s\nДобавлена пользователем с ID: %d",
		expense.Sum,
		expense.Description,
		expense.Email,
		expense.CreationDate.Format("2006-01-02 15:04:05"),
		// userName,
		expense.UserID)
	err := bot.sendMsg(chatID, expenseMsg)
	return err
}
