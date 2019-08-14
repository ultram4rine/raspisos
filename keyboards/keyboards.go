package keyboards

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ultram4rine/raspisos/parsing"
)

func CreateMainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	e1, _ := strconv.ParseInt(strings.TrimPrefix("\\U0001F4DA", "\\U"), 16, 32)
	e2, _ := strconv.ParseInt(strings.TrimPrefix("\\U0001F4DD", "\\U"), 16, 32)

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(string(e1)+" Занятия"),
			tgbotapi.NewKeyboardButton(string(e2)+" Сессия"),
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
