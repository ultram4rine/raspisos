package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ultram4rine/raspisos/parsing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func handler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Raspisos here"))
}

//Config struct contains bot's token
type Config struct {
	TgBotToken string `json:"TelegramBotToken"`
}

var conf Config

func makeConfig(path string) error {
	conffile, err := os.Open(path)
	if err != nil {
		return err
	}

	confdata, err := ioutil.ReadAll(conffile)
	if err != nil {
		return err
	}
	err = conffile.Close()

	err = json.Unmarshal(confdata, &conf)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var (
		/*facultiess = [...]string{
			"bf", "gf", "gl",
			"idpo", "ii", "imo",
			"ifk", "ifg", "ih",
			"mm", "sf", "fi", "knt",
			"fn", "fnp", "fps", "fppso",
			"ff", "fp", "ef",
			"uf", "kgl", "cre",
		}*/
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

		faculty     = "mm"
		group       = "211"
		schedule    = "lesson"
		address     = "https://www.sgu.ru/schedule/" + faculty + "/do/" + group + "/" + schedule
		xmlschedule parsing.XMLStruct
	)

	err := makeConfig("conf.json")
	if err != nil {
		log.Fatal("Error with making config:", err)
	}

	bot, err := tgbotapi.NewBotAPI(conf.TgBotToken)
	if err != nil {
		log.Fatal("Error with creaction of the bot: ", err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook("/" + bot.Token)

	http.HandleFunc("/", handler)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	//Bot answers
	for update := range updates {
		if update.Message == nil {
			continue
		} else {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			day := update.Message.Command()
			switch day {
			case "start":
				msg.Text = "Let's get it"
			default:
				msg.Text, err = MakeMessage(schedule+"_"+faculty+"_"+group+".xml", address, day, xmlschedule)
				if err != nil {
					log.Println(err)
				}
			}
			msg.ParseMode = "markdown"
			bot.Send(msg)
		}
	}
}

//MakeMessage make message to answer with schedule
func MakeMessage(filepath, address, day string, xmlschedule parsing.XMLStruct) (msgtext string, err error) {
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
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
