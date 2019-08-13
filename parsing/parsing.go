package parsing

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	abbr "github.com/ultram4rine/raspisos/abbreviation"
)

//XMLStruct is a struct for parsing XML file from https://sgu.ru/schedule/...
//This struct was generated on https://onlinetool.io/xmltogo, unnecessary fields was deleted
type XMLStruct struct {
	Worksheet []struct {
		Table struct {
			Row []struct {
				Cell []struct {
					Data struct {
						Text string `xml:",chardata"`
					} `xml:"Data"`
				} `xml:"Cell"`
			} `xml:"Row"`
		} `xml:"Table"`
	} `xml:"Worksheet"`
}

//Fac is a struct describing faculty
type Fac struct {
	Name string
	Link string
	Abbr string
}

//Lesson is a struct describing a lesson
type Lesson struct {
	Time         string
	TypeofWeek   string
	TypeofLesson string
	Name         string
	SubGroup     string
	Teacher      string
	Classroom    string
	Data         string
}

//Day is a type describing a day
type Day struct {
	Lessons [10]Lesson
}

//Week is a type describing week
type Week struct {
	Days [7]Day
}

func MakeWeek(xmlschedule XMLStruct) (week Week) {
	rowsize := len(xmlschedule.Worksheet[0].Table.Row)
	for i := 2; i < rowsize; i++ {
		cellsize := len(xmlschedule.Worksheet[0].Table.Row[i].Cell)
		for j := 0; j < cellsize; j++ {
			data := xmlschedule.Worksheet[0].Table.Row[i].Cell[j].Data.Text
			week.Days[j].Lessons[i].Data = data
		}
	}
	return week
}

func (week Week) makeTime() (d Day) {
	return week.Days[0]
}

func (week Week) makeDay(day string) (d Day) {
	switch day {
	case "monday":
		d = week.Days[1]
	case "tuesday":
		d = week.Days[2]
	case "wednesday":
		d = week.Days[3]
	case "thursday":
		d = week.Days[4]
	case "friday":
		d = week.Days[5]
	case "saturday":
		d = week.Days[6]
	}

	return d
}

func formatLesson(s string) (l Lesson) {
	var (
		n = strings.Count(s, "\n")
		k int
	)

	if strings.Contains(s, "знам.") {
		l.TypeofWeek = "Знаменатель"
	}
	if strings.Contains(s, "чис.") {
		l.TypeofWeek = "Числитель"
	}
	if strings.Contains(s, "лек.") {
		l.TypeofLesson = "Лекция"
	}
	if strings.Contains(s, "пр.") {
		l.TypeofLesson = "Практика"
	}

	k = strings.Index(s, ".")
	if k != -1 {
		s = s[k+1:]
		k = strings.IndexRune(s, '\n')
		l.Name = s[:k]
		s = s[k+1:]

		if n%5 == 0 {
			k = strings.IndexRune(s, '\n')
			sub := s[:k]
			if strings.Contains(sub, "под.") {
				if strings.Contains(sub, "1") {
					l.SubGroup = "1"
				} else if strings.Contains(sub, "2") {
					l.SubGroup = "2"
				}
			}
			s = s[k+1:]
		}

		k = strings.IndexRune(s, '\n')
		l.Teacher = s[:k]
		s = s[k+1:]

		k = strings.IndexRune(s, '\n')
		l.Classroom = s[:k]
		s = s[k+1:]
	}

	return l
}

func MakeLesson(week Week, day string, number int) (l Lesson) {
	var (
		t = week.makeTime()
		d = week.makeDay(day)
	)
	switch number {
	case 1:
		l.Data = d.Lessons[2].Data
		l = formatLesson(l.Data)
		l.Time = strings.Replace(t.Lessons[2].Data, "- ", "-", -1)
	case 2:
		l.Data = d.Lessons[3].Data
		l = formatLesson(l.Data)
		l.Time = strings.Replace(t.Lessons[3].Data, "- ", "-", -1)
	case 3:
		l.Data = d.Lessons[4].Data
		l = formatLesson(l.Data)
		l.Time = strings.Replace(t.Lessons[4].Data, "- ", "-", -1)
	case 4:
		l.Data = d.Lessons[5].Data
		l = formatLesson(l.Data)
		l.Time = strings.Replace(t.Lessons[5].Data, "- ", "-", -1)
	case 5:
		l.Data = d.Lessons[6].Data
		l = formatLesson(l.Data)
		l.Time = strings.Replace(t.Lessons[6].Data, "- ", "-", -1)
	case 6:
		l.Data = d.Lessons[7].Data
		l = formatLesson(l.Data)
		l.Time = strings.Replace(t.Lessons[7].Data, "- ", "-", -1)
	case 7:
		l.Data = d.Lessons[8].Data
		l = formatLesson(l.Data)
		l.Time = strings.Replace(t.Lessons[8].Data, "- ", "-", -1)
	case 8:
		l.Data = d.Lessons[9].Data
		l = formatLesson(l.Data)
		l.Time = strings.Replace(t.Lessons[9].Data, "- ", "-", -1)
	}

	return l
}

//GetFacs gets all faculties
func GetFacs() ([]Fac, error) {
	var (
		addr      = "https://www.sgu.ru/schedule"
		fac       Fac
		faculties []Fac
	)

	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(".panes_item__type_group ul li").Each(func(i int, s *goquery.Selection) {
		fac.Name = s.Find("a").Text()
		href, _ := s.Find("a").Attr("href")
		fac.Link = strings.Split(href, "/")[2]
		fac.Abbr = abbr.Faculty(fac.Name)
		faculties = append(faculties, fac)
	})

	return faculties, nil
}
