package emoji

import (
	"strconv"
	"strings"
)

//Emoji type describe emoji
type Emoji struct {
	Name string
	Code string
}

var emojiList = []Emoji{{"bell", "\\U0001F514"}, {"books", "\\U0001F4DA"}, {"gear", "\\U00002699"}, {"memo", "\\U0001F4DD"}}

//Parse function parse all used emoji
func Parse() (map[string]string, error) {
	emojiMap := make(map[string]string)

	for _, e := range emojiList {
		eInt, err := strconv.ParseInt(strings.TrimPrefix(e.Code, "\\U"), 16, 32)
		if err != nil {
			return nil, err
		}

		emojiMap[e.Name] = string(eInt)
	}

	return emojiMap, nil
}
