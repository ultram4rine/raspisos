package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ultram4rine/raspisos/config"
	"github.com/ultram4rine/raspisos/data"
	"github.com/ultram4rine/raspisos/emoji"
	"github.com/ultram4rine/raspisos/keyboards"
	"github.com/ultram4rine/raspisos/schedule"
	"github.com/ultram4rine/raspisos/www"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	var (
		confPath = "conf.json"
		days     = []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"}
		newUser  data.User
		users    = make(map[int64]data.User)
		//TODO: Save faculties and groups of user to DB or json files
		//TODO: webDesign and translator option
		//webDesign  bool
		//translator bool
		//TODO: subGroup numbers
		//subGroupNum
		//englishGroupNum
		//translatorGroupNum
		//TODO: selection of group and faculty
		//TODO: schedule type "session"
		//FIXME: non full output for multiply lessons in one time

		faculty     string
		group       string
		xmlschedule schedule.XMLStruct
	)

	err := config.ParseConfig(confPath)
	if err != nil {
		log.Fatal("Error parsing config: ", err)
	}

	emojiMap, err := emoji.Parse()
	if err != nil {
		log.Fatal("Error parsing emoji: ", err)
	}

	faculties, err := www.GetFacs()
	if err != nil {
		log.Fatal("Error getting faculties:", err)
	}
	go func() {
		for range time.Tick(time.Hour * 24) {
			faculties, err = www.GetFacs()
			if err != nil {
				log.Println("Error updating faculties:", err)
			}
		}
	}()

	_, err = time.LoadLocation("Europe/Saratov")
	if err != nil {
		log.Fatal("Error getting time zone: ", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.Conf.TgBotToken)
	if err != nil {
		log.Fatal("Error with creaction of the bot: ", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("Error getting updates channel of the bot: ", err)
	}

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			cmd := update.Message.Command()
			if cmd != "" {
				switch cmd {
				case "start":
					var keyboard tgbotapi.ReplyKeyboardMarkup

					if user, ok := users[update.Message.Chat.ID]; ok {
						msg.Text = "I already have your settings.\nФакультет: " + user.Faculty + "\nГруппа: " + user.Group
						keyboard = keyboards.CreateMainKeyboard(emojiMap)
					} else {
						users[update.Message.Chat.ID] = newUser

						msg.Text = "New user! Set your group."
						keyboard = keyboards.CreateFirstKeyboard(emojiMap)
					}

					msg.ReplyMarkup = keyboard
				case "help":
					msg.Text = "Help message"
				case "src":
					msg.Text = "Код приложения доступен по ссылке.\nhttps://github.com/ultram4rine/raspisos"
				default:
					msg.Text = "Unknown command, type \"/help\" for help"
				}
			} else {
				text := update.Message.Text

				switch text {
				case emojiMap["books"] + " Занятия":
					keyboard := keyboards.CreateDaysKeyboard()

					msg.ReplyMarkup = keyboard
					msg.Text = "Выберите день недели"
				case "Изменить группу", "Задать группу":
					keyboard := keyboards.CreateFacsKeyboard(faculties)

					msg.ReplyMarkup = keyboard
					msg.Text = "Выберите факультет"
				}
			}

			msg.ParseMode = "markdown"
			bot.Send(msg)
		}
		if update.CallbackQuery != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

			ans := update.CallbackQuery.Data
			typo := strings.Split(ans, "&")[1]
			switch typo {
			case "fac":
				newUser.Faculty = strings.Split(ans, "&")[0]

				educationTypes, _ := www.GetTypesofEducation(strings.Split(ans, "&")[0])

				keyboard := keyboards.CreateGroupsKeyboard(strings.Split(ans, "&")[0], educationTypes[0], educationTypes)

				msgedit := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Образование: *"+strings.ToLower(educationTypes[0])+"*. Выберите группу")
				msgedit.ReplyMarkup = keyboard
				msgedit.ParseMode = "markdown"
				bot.Send(msgedit)
			case "group":
				newUser.Group = strings.Split(ans, "&")[0]

				users[update.CallbackQuery.Message.Chat.ID] = newUser

				msgdelete := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				bot.Send(msgdelete)

				keyboard := keyboards.CreateMainKeyboard(emojiMap)
				msg.ReplyMarkup = keyboard
				msg.Text = "Done!"
				bot.Send(msg)
			case "day":
				if user, ok := users[update.CallbackQuery.Message.Chat.ID]; ok {
					day := strings.Split(ans, "&")[0]
					address := "https://www.sgu.ru/schedule/" + user.Faculty + "/do/" + user.Group + "/lesson"

					if contains(days, day) {
						msg.Text, err = makeLessonMsg("lesson"+"_"+faculty+"_"+group+".xml", address, day, xmlschedule)
						if err != nil {
							log.Println(err)
						}

						msg.ParseMode = "markdown"
						bot.Send(msg)
					} else {
						if day == "Сегодня" {

						} else if day == "Завтра" {

						}
					}
				} else {
					keyboard := keyboards.CreateFirstKeyboard(emojiMap)

					msg.ReplyMarkup = keyboard
					msg.Text = "Make your schedule"
					bot.Send(msg)
				}
			}
		}
	}
}

func contains(slice []string, x string) bool {
	for _, v := range slice {
		if x == v {
			return true
		}
	}

	return false
}

//makeLessonMsg make message to answer with schedule
func makeLessonMsg(filepath, address, day string, xmlschedule schedule.XMLStruct) (msgtext string, err error) {
	err = www.DownloadFileFromURL(filepath, address)
	if err != nil {
		return "", err
	}

	xmlfile, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	xmldata, err := ioutil.ReadAll(xmlfile)
	if err != nil {
		return "", err
	}

	err = xmlfile.Close()
	if err != nil {
		return "", err
	}

	err = xml.Unmarshal(xmldata, &xmlschedule)
	if err != nil {
		return "", err
	}

	week := schedule.MakeWeek(xmlschedule)

	for i := 1; i < 9; i++ {
		l := schedule.MakeLesson(week, day, i)
		msgtext += "*" + l.Time + "*\n"

		if l.Name != "" {
			if l.TypeofWeek != "" {
				msgtext += "*Неделя:* " + l.TypeofWeek + "\n"
			}

			if l.TypeofLesson != "" {
				msgtext += "*Занятие:* " + l.TypeofLesson + "\n"
			}

			msgtext += "*Предмет:* " + l.Name + "\n"

			if l.SubGroup != "" {
				msgtext += "*Подгруппа:* " + l.SubGroup + "\n"
			}

			msgtext += "*Преподаватель:* " + l.Teacher + "\n"

			msgtext += "*Аудитория:* " + l.Classroom + "\n\n"
		} else {
			msgtext += "_Пары нет_\n\n"
		}
	}

	err = os.Remove(filepath)
	if err != nil {
		return msgtext, err
	}

	return msgtext, nil
}
