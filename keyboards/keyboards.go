package keyboards

import (
	"github.com/ultram4rine/raspisos/data"
	"github.com/ultram4rine/raspisos/www"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CreateMainKeyboard(emojiMap map[string]string) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(emojiMap["books"]+" Занятия"),
			tgbotapi.NewKeyboardButton(emojiMap["memo"]+" Сессия"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить группу"),
			tgbotapi.NewKeyboardButton(emojiMap["gear"]+" Настройки"),
		),
	)

	return keyboard
}

func CreateFacsKeyboard(faculties []data.Fac) tgbotapi.InlineKeyboardMarkup {
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

func CreateGroupsKeyboard(facLink, educationType string, educationTypes []string) *tgbotapi.InlineKeyboardMarkup {
	var (
		btn  tgbotapi.InlineKeyboardButton
		row  []tgbotapi.InlineKeyboardButton
		rows [][]tgbotapi.InlineKeyboardButton
	)

	groups, _ := www.GetGroups(facLink, educationType)

	for i, g := range groups {
		if i%3 != 0 || i == 0 {
			btn = tgbotapi.NewInlineKeyboardButtonData(g, g+"&group")
			row = append(row, btn)
		} else {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
			btn = tgbotapi.NewInlineKeyboardButtonData(g, g+"&group")
			row = append(row, btn)
		}
	}
	rows = append(rows, row)

	row = []tgbotapi.InlineKeyboardButton{}
	for i, t := range educationTypes {
		if t == educationType {
			continue
		} else {
			if i%3 != 0 || i == 0 {
				btn = tgbotapi.NewInlineKeyboardButtonData(t, t+"&education")
				row = append(row, btn)
			} else {
				rows = append(rows, row)
				row = []tgbotapi.InlineKeyboardButton{}
				btn = tgbotapi.NewInlineKeyboardButtonData(t, t+"&education")
				row = append(row, btn)
			}
		}
	}
	rows = append(rows, row)

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	keyboard.InlineKeyboard = rows

	return &keyboard
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
