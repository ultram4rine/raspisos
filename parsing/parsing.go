package parsing

func runeIndex(text []rune, substr string) int {
	subrunes := []rune(substr)
	for i := range text {
		found := true
		for j := range subrunes {
			if text[i+j] != subrunes[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}

func strSub(str string, i, j int) (substr string) {
	runes := []rune(str)
	substr = string(runes[i:j])
	return substr
}

func strErase(str string, i, j int) (strerased string) {
	runes := []rune(str)
	strerased = string(append(runes[:i], runes[j:]...))
	return strerased
}

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
	Data string
}

//MakeLesson makes lesson object from XML struct
func MakeLesson(xmlschedule XMLStruct) [][]Lesson {
	rowsize := len(xmlschedule.Worksheet[0].Table.Row)
	lesson := make([][]Lesson, rowsize)
	for i := 2; i < rowsize; i++ {
		cellsize := len(xmlschedule.Worksheet[0].Table.Row[i].Cell)
		lesson[i] = make([]Lesson, cellsize)
	}
	for i := 2; i < rowsize; i++ {
		cellsize := len(xmlschedule.Worksheet[0].Table.Row[i].Cell)
		for j := 0; j < cellsize; j++ {
			data := xmlschedule.Worksheet[0].Table.Row[i].Cell[j].Data.Text
			lesson[i][j].Data = data
		}
	}
	return lesson
}
