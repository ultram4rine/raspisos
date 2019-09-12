package config

import (
	"io/ioutil"
	"os"

	jsoniter "github.com/json-iterator/go"
)

//Conf contains main data
var Conf struct {
	TgBotToken string `json:"TelegramBotToken"`
}

//ParseConfig to parse config
func ParseConfig(confPath string) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	confFile, err := os.Open(confPath)
	if err != nil {
		return err
	}
	defer confFile.Close()

	confData, err := ioutil.ReadAll(confFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(confData, &Conf)
	if err != nil {
		return err
	}

	return nil
}
