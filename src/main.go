package main

import (
	"ad-spam-tgbot/db"
	// "bufio"
	// "fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/exec"
	"strconv"
	// "strings"
	"time"
)

var mainMenuKeysStart = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Набрать сообщение", "create"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить получателя", "add_reciever"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Настроить интервал", "set_interval"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Запустить", "start"),
	),
)

var mainMenuKeysStop = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Набрать сообщение", "create"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить получателя", "add_reciever"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Настроить интервал", "set_interval"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Остановить", "stop"),
	),
)

func CreateAdMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msgText := db.Message{ID: 1, Text: update.Message.Text}
	db.DB.Save(msgText)
	msg := tgbotapi.NewMessage(update.Message.From.ID, "Сообщение сохранено")
	bot.Send(msg)
}

func AddNewReciever(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	new := update.Message.Text
	newReciever := db.Reciever{Nickname: new}
	db.DB.Create(&newReciever)
	msg := tgbotapi.NewMessage(update.Message.From.ID, "Получатель добавлен")
	bot.Send(msg)

}

func SetTimeInterval(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	newInterval, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.From.ID, "Ошибка!")
		bot.Send(msg)
		return
	}
	interval = newInterval
	msg := tgbotapi.NewMessage(update.Message.From.ID, "Интервал установлен")
	bot.Send(msg)
}

func Sender() {
	ticker := time.NewTicker(time.Minute * time.Duration(interval))
	for senderState {
		select {
		case <-ticker.C:
			cmd := exec.Command("/Users/yalu/Developer/GO_projects/ad-spam-tgbot/src/py_sender/tg.py")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			log.Println(cmd.Run())
		default:
		}
	}
}

var (
	senderState bool
	interval    int = 20
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("SPAMBOT_APITOKEN"))
	if err != nil {
		panic(err)
	}

	userStates := make(map[int64]string)

	db.ConnectDB()

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			go UpdateMessageHandler(bot, update, userStates)
		} else if update.CallbackQuery != nil {
			go UpdateCallbackHandler(bot, update, userStates)
		}
	}

}

func UpdateMessageHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update, userStates map[int64]string) {
	if update.Message.IsCommand() && update.Message.Command() == "start" {
		msg := tgbotapi.NewMessage(update.Message.From.ID, ":)")

		if senderState {
			msg.ReplyMarkup = mainMenuKeysStop
		} else {
			msg.ReplyMarkup = mainMenuKeysStart
		}

		bot.Send(msg)
	} else {
		switch userStates[update.Message.From.ID] {
		case "waiting_for_new_msg":
			CreateAdMessage(bot, update)

		case "waiting_for_new_reciever":
			AddNewReciever(bot, update)

		case "set_interval":
			SetTimeInterval(bot, update)

		}
		userStates[update.Message.From.ID] = ""
	}

}

func UpdateCallbackHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update, userStates map[int64]string) {
	userID := update.CallbackQuery.From.ID
	callerID := update.CallbackQuery.Message.MessageID
	reMsg := tgbotapi.NewEditMessageTextAndMarkup(userID, callerID, update.CallbackQuery.Message.Text, mainMenuKeysStop)
	switch update.CallbackData() {
	case "create":
		msg := tgbotapi.NewMessage(userID, "Наберите рекламное сообщение:")
		userStates[userID] = "waiting_for_new_msg"
		bot.Send(msg)
	case "add_reciever":
		msg := tgbotapi.NewMessage(userID, "Отправьте ID нового получателя:")
		userStates[userID] = "waiting_for_new_reciever"
		bot.Send(msg)
	case "start":
		senderState = true
		go Sender()
		bot.Request(reMsg)
	case "set_interval":
		msg := tgbotapi.NewMessage(userID, "Укажите временной интервал между отправкой сообщений в минутах:")
		userStates[userID] = "set_interval"
		bot.Request(msg)
	case "stop":
		senderState = false
		reMsg.ReplyMarkup = &mainMenuKeysStart
		bot.Request(reMsg)
	case "5min":
	}
}
