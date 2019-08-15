package emoji

import (
	"log"
	"strconv"
	"strings"
)

const (
	books = "\\U0001F4DA"
	memo  = "\\U0001F4DD"
	bell  = "\\U0001F514"
)

func Parse() (e1, e2, e3 string, err error) {
	e1int, err := strconv.ParseInt(strings.TrimPrefix(books, "\\U"), 16, 32)
	if err != nil {
		log.Printf("Error parsing books emoji: %s", err)
		return "", "", "", err
	}

	e2int, err := strconv.ParseInt(strings.TrimPrefix(memo, "\\U"), 16, 32)
	if err != nil {
		log.Printf("Error parsing memo emoji: %s", err)
		return "", "", "", err
	}

	e3int, err := strconv.ParseInt(strings.TrimPrefix(bell, "\\U"), 16, 32)
	if err != nil {
		log.Printf("Error parsing bell emoji: %s", err)
		return "", "", "", err
	}

	return string(e1int), string(e2int), string(e3int), nil
}
