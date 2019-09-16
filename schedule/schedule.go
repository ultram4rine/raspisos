package schedule

import (
	"strings"
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

//Lesson is a struct describing a lesson
type Lesson struct {
	Time         string
	TypeofWeek   string
	TypeofLesson string
	Name         string
	SubGroup     string
	Teacher      string
	Classroom    string
	Special      string
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
	case "Понедельник":
		d = week.Days[1]
	case "Вторник":
		d = week.Days[2]
	case "Среда":
		d = week.Days[3]
	case "Четверг":
		d = week.Days[4]
	case "Пятница":
		d = week.Days[5]
	case "Суббота":
		d = week.Days[6]
	}

	return d
}

func formatLesson(s string) []Lesson {
	var k, length int

	lStrings := strings.Split(s, "\n\n")

	if len(lStrings) > 2 {
		for i, lesson := range lStrings {
			if !strings.Contains(lesson, "лек.") && !strings.Contains(lesson, "пр.") && lesson != " " {
				if i == len(lStrings)-1 {
					lStrings[i-1] += "\n" + lStrings[i]
					lStrings = lStrings[:len(lStrings)-1]
				} else {
					lStrings[i-1] += "\n" + lStrings[i]
					lStrings = append(lStrings[:i], lStrings[i+1:]...)
				}
			}
		}
	}

	switch len(lStrings) == 1 {
	case true:
		length = 1
	case false:
		length = len(lStrings) - 1
	}

	lessons := make([]Lesson, length)

	for i, lesson := range lStrings {
		if lesson == " " && i != 0 {
			continue
		}

		lesson += "\n"
		n := strings.Count(lesson, "\n")

		if strings.Contains(lesson, "знам.") {
			lessons[i].TypeofWeek = "Знаменатель"
			k = strings.Index(lesson, ".")
			lesson = lesson[k+1:]
		}
		if strings.Contains(lesson, "чис.") {
			lessons[i].TypeofWeek = "Числитель"
			k = strings.Index(lesson, ".")
			lesson = lesson[k+1:]
		}
		if strings.Contains(lesson, "лек.") {
			lessons[i].TypeofLesson = "Лекция"
			k = strings.Index(lesson, ".")
			lesson = lesson[k+1:]
		}
		if strings.Contains(lesson, "пр.") {
			lessons[i].TypeofLesson = "Практика"
			k = strings.Index(lesson, ".")
			lesson = lesson[k+1:]
		}

		k = strings.Index(lesson, ".")
		if k != -1 {
			k = strings.IndexRune(lesson, '\n')
			lessons[i].Name = lesson[:k]
			lessons[i].Name = strings.TrimSpace(lessons[i].Name)
			lesson = lesson[k+1:]

			if n%4 == 0 {
				k = strings.IndexRune(lesson, '\n')
				sub := lesson[:k]
				if strings.Contains(sub, "под.") {
					if strings.Contains(sub, "1") {
						lessons[i].SubGroup = "1"
					} else if strings.Contains(sub, "2") {
						lessons[i].SubGroup = "2"
					}
				} else {
					lessons[i].Special = sub
				}
				lesson = lesson[k+1:]
			}

			k = strings.IndexRune(lesson, '\n')
			if k != -1 && !strings.ContainsAny(lesson[:k], "0123456789") {
				lessons[i].Teacher = lesson[:k]
				lesson = lesson[k+1:]
			}

			k = strings.IndexRune(lesson, '\n')
			if k != -1 {
				lessons[i].Classroom = lesson[:k]
				lesson = lesson[k+1:]
			}
		}
	}

	return lessons
}
func formatAndsetTime(number int, t, d Day, lessons []Lesson) []Lesson {
	lessons = formatLesson(d.Lessons[number].Data)
	for i := range lessons {
		lessons[i].Time = strings.Replace(t.Lessons[number].Data, "- ", " - ", -1)
	}

	return lessons
}

func MakeLesson(week Week, day string, number int) (lessons []Lesson) {
	var (
		t = week.makeTime()
		d = week.makeDay(day)
	)

	switch number {
	case 1:
		lessons = formatAndsetTime(2, t, d, lessons)
	case 2:
		lessons = formatAndsetTime(3, t, d, lessons)
	case 3:
		lessons = formatAndsetTime(4, t, d, lessons)
	case 4:
		lessons = formatAndsetTime(5, t, d, lessons)
	case 5:
		lessons = formatAndsetTime(6, t, d, lessons)
	case 6:
		lessons = formatAndsetTime(7, t, d, lessons)
	case 7:
		lessons = formatAndsetTime(8, t, d, lessons)
	case 8:
		lessons = formatAndsetTime(9, t, d, lessons)
	}

	return lessons
}
