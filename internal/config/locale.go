package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func SetupLocale() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English) // Default language
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err := bundle.LoadMessageFile("resources/en/message.en.toml")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = bundle.LoadMessageFile("resources/ru/message.ru.toml")
	if err != nil {
		fmt.Println(err.Error())
	}
	return bundle
}
