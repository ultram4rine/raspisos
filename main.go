package main

import (
	"database/sql"
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
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func main() {
	var (
		confPath  = "conf.json"
		days      = []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"}
		newUser   data.User
		usersList []data.User
		usersMap  = make(map[int64]data.User)
		//TODO: webDesign and translator option
		//webDesign  bool
		//translator bool
		//TODO: subGroup numbers
		//subGroupNum
		//englishGroupNum
		//translatorGroupNum
		//TODO: selection of group and faculty
		//TODO: schedule type "session"

		xmlschedule schedule.XMLStruct
	)

	err := config.ParseConfig(confPath)
	if err != nil {
		log.Fatal("Error parsing config: ", err)
	}

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal("Error getting users list: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user data.User

		if err := rows.Scan(&user.ID, &user.Faculty, &user.Group, &user.Notifications, pq.Array(&user.IgnoreList)); err != nil {
			log.Fatal("Error scanning user: ", err)
		}

		usersList = append(usersList, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Rows error: ", err)
	}

	for _, user := range usersList {
		usersMap[user.ID] = user
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

					if user, ok := usersMap[update.Message.Chat.ID]; ok {
						msg.Text = "I already have your settings.\n Факультет: " + user.Faculty + "\n Группа: " + user.Group
						keyboard = keyboards.CreateMainKeyboard(emojiMap)
					} else {
						usersMap[update.Message.Chat.ID] = newUser

						msg.Text = "New user! Set your group."
						keyboard = keyboards.CreateFirstKeyboard(emojiMap)
					}

					msg.ReplyMarkup = keyboard
				case "help":
					msg.Text = "Help message"
				case "src":
					msg.Text = "Код приложения доступен по ссылке.\n https://github.com/ultram4rine/raspisos"
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
				case emojiMap["memo"] + " Сессия":

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
				newUser.ID = update.CallbackQuery.Message.Chat.ID
				newUser.Notifications = false

				usersMap[update.CallbackQuery.Message.Chat.ID] = newUser

				msgdelete := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				bot.Send(msgdelete)

				var userID int64

				err = db.QueryRow("SELECT id FROM users WHERE id = $1", update.CallbackQuery.Message.Chat.ID).Scan(&userID)
				if err != nil {
					if err != sql.ErrNoRows {
						log.Println("Error checking user settings existing in DB: ", err)
					} else {
						_, err := db.Exec("INSERT INTO users (id, fac, groupnum, notifications, ignorelist) VALUES ($1, $2, $3, $4, $5)", newUser.ID, newUser.Faculty, newUser.Group, newUser.Notifications, pq.Array(newUser.IgnoreList))
						if err != nil {
							log.Println("Error inserting user settings to db: ", err)

							msg.Text = "Error inserting your settings to db"
							bot.Send(msg)
						} else {
							keyboard := keyboards.CreateMainKeyboard(emojiMap)
							msg.ReplyMarkup = keyboard
							msg.Text = "Done!"
							bot.Send(msg)
						}
					}
				} else {
					_, err := db.Exec("UPDATE users SET (id, fac, groupnum, notifications, ignorelist) = ($1, $2, $3, $4, $5)", newUser.ID, newUser.Faculty, newUser.Group, newUser.Notifications, pq.Array(newUser.IgnoreList))
					if err != nil {
						log.Println("Error updating user settings to db: ", err)

						msg.Text = "Error updating your settings to db"
						bot.Send(msg)
					} else {
						keyboard := keyboards.CreateMainKeyboard(emojiMap)
						msg.ReplyMarkup = keyboard
						msg.Text = "Done!"
						bot.Send(msg)
					}
				}
			case "day":
				if user, ok := usersMap[update.CallbackQuery.Message.Chat.ID]; ok {
					day := strings.Split(ans, "&")[0]
					address := "https://www.sgu.ru/schedule/" + user.Faculty + "/do/" + user.Group + "/lesson"

					if contains(days, day) {
						msg.Text, err = makeLessonMsg(address, day, xmlschedule)
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
func makeLessonMsg(address, day string, xmlschedule schedule.XMLStruct) (msgtext string, err error) {
	data, err := www.DownloadFileFromURL(address)
	if err != nil {
		return "", err
	}
	defer data.Close()

	xmldata, err := ioutil.ReadAll(data)
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
		msgtext += "*" + l[0].Time + "*\n "

		for j := range l {

			if l[j].Name != "" {
				if l[j].TypeofWeek != "" {
					msgtext += "*Неделя:* " + l[j].TypeofWeek + "\n "
				}

				if l[j].TypeofLesson != "" {
					msgtext += "*Занятие:* " + l[j].TypeofLesson + "\n "
				}

				if l[j].Special != "" {
					msgtext += "*Доп. курс:* " + l[j].Special + "\n "
				}

				msgtext += "*Предмет:* " + l[j].Name + "\n "

				if l[j].SubGroup != "" {
					msgtext += "*Подгруппа:* " + l[j].SubGroup + "\n "
				}

				msgtext += "*Преподаватель:* " + l[j].Teacher + "\n "

				msgtext += "*Аудитория:* " + l[j].Classroom + "\n \n "
			} else {
				msgtext += "_Пары нет_\n \n "
			}
		}
	}

	return msgtext, nil
}
