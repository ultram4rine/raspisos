package keyboards

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ultram4rine/raspisos/emoji"
	"github.com/ultram4rine/raspisos/parsing"
)

func CreateMainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	books, memo, bell, err := emoji.Parse()
	if err != nil {
		log.Printf("Error parsing emojies: %s", err)
	}

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(books+" Занятия"),
			tgbotapi.NewKeyboardButton(memo+" Сессия"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bell+"Уведомления"),
		),
	)

	return keyboard
}

func CreateFacsKeyboard(faculties []parsing.Fac) tgbotapi.InlineKeyboardMarkup {
	var (
		btn  tgbotapi.InlineKeyboardButton
		row  []tgbotapi.InlineKeyboardButton
		rows [][]tgbotapi.InlineKeyboardButton
	)

	for i, f := range faculties {
		if i%3 != 0 || i == 0 {
			btn = tgbotapi.NewInlineKeyboardButtonData(f.Abbr, f.Link+"&fac")
			row = append(row, btn)
		} else {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
			btn = tgbotapi.NewInlineKeyboardButtonData(f.Abbr, f.Link+"&fac")
			row = append(row, btn)
		}
	}
	rows = append(rows, row)

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	keyboard.InlineKeyboard = rows

	return keyboard
}

func CreateDaysKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "Сегодня&day"),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "Завтра&day"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Понедельник", "Понедельник&day"),
			tgbotapi.NewInlineKeyboardButtonData("Вторник", "Вторник&day"),
			tgbotapi.NewInlineKeyboardButtonData("Среда", "Среда&day"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Четверг", "Четверг&day"),
			tgbotapi.NewInlineKeyboardButtonData("Пятница", "Пятница&day"),
			tgbotapi.NewInlineKeyboardButtonData("Суббота", "Суббота&day"),
		),
	)

	return keyboard
}
