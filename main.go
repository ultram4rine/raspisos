package main

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ultram4rine/raspisos/keyboards"
	"github.com/ultram4rine/raspisos/parsing"

	"github.com/json-iterator/go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func handler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Raspisos here"))
}

var conf struct {
	TgBotToken string `json:"TelegramBotToken"`
}

func main() {
	var (
		confpath = "conf.json"
		days     = []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"}
		userMap  = make(map[int]string)
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
		schedule    = "lesson"
		xmlschedule parsing.XMLStruct
	)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	faculties, err := parsing.GetFacs()
	if err != nil {
		log.Println("Error getting faculties:", err)
	}

	confFile, err := os.Open(confpath)
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}

	confData, err := ioutil.ReadAll(confFile)
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}
	err = confFile.Close()

	err = json.Unmarshal(confData, &conf)
	if err != nil {
		log.Fatal("Error unmarshalling config:", err)
	}

	bot, err := tgbotapi.NewBotAPI(conf.TgBotToken)
	if err != nil {
		log.Fatal("Error with creaction of the bot: ", err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook("/" + bot.Token)

	_, err = time.LoadLocation("Europe/Saratov")
	if err != nil {
		log.Printf("Error getting time zone: %s", err)
	}

	http.HandleFunc("/", handler)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	//Bot answers
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			cmd := update.Message.Command()
			if cmd != "" {
				switch cmd {
				case "start":
					keyboard := keyboards.CreateMainKeyboard()

					msg.ReplyMarkup = keyboard
					msg.Text = "Started"
				case "help":
					msg.Text = "Help message"
				default:
					msg.Text = "Unknown command, type \"/help\" for help"
				}
			} else {
				text := update.Message.Text

				e1, _ := strconv.ParseInt(strings.TrimPrefix("\\U0001F4DA", "\\U"), 16, 32)

				switch text {
				case string(e1) + " Занятия":
					keyboard := keyboards.CreateDaysKeyboard()

					msg.ReplyMarkup = keyboard
					msg.Text = "Выберите день недели"
				case "Изменить группу":
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
				userMap[update.CallbackQuery.From.ID] = strings.Split(ans, "&")[0]
				log.Println(userMap[update.CallbackQuery.From.ID])
				msgedit := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Выберите группу")
				msgedit.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("211", "211&group"),
					),
				)).ReplyMarkup
				bot.Send(msgedit)
			case "group":
				userMap[update.CallbackQuery.From.ID] += "&" + strings.Split(ans, "&")[0]
				log.Println(userMap[update.CallbackQuery.From.ID])
				msgedit := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Done")
				bot.Send(msgedit)
			case "day":
				day := strings.Split(ans, "&")[0]
				if contains(days, day) {
					faculty = strings.Split(userMap[update.CallbackQuery.From.ID], "&")[0]
					group = strings.Split(userMap[update.CallbackQuery.From.ID], "&")[1]
					address := "https://www.sgu.ru/schedule/" + faculty + "/do/" + group + "/" + schedule

					msg.Text, err = makeLessonMsg(schedule+"_"+faculty+"_"+group+".xml", address, day, xmlschedule)
					if err != nil {
						log.Println(err)
					}

					msg.ParseMode = "markdown"
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
func makeLessonMsg(filepath, address, day string, xmlschedule parsing.XMLStruct) (msgtext string, err error) {
	err = DownloadFileFromURL(filepath, address)
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

	week := parsing.MakeWeek(xmlschedule)

	for i := 1; i < 9; i++ {
		l := parsing.MakeLesson(week, day, i)
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

//DownloadFileFromURL for download
func DownloadFileFromURL(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("bad status: " + resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
