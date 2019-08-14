package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ultram4rine/raspisos/parsing"
)

func CreateMainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("U+1F4DA Занятия"),
			tgbotapi.NewKeyboardButton("U+1F4DD Сессия"),
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
