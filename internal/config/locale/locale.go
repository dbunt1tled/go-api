package locale

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	localeBundleInstance *i18n.Bundle    //nolint:gochecknoglobals // singleton
	localizerInstance    *i18n.Localizer //nolint:gochecknoglobals // singleton
	m                    sync.Mutex
)

func GetLocaleBundleInstance() *i18n.Bundle {
	if localeBundleInstance == nil {
		localeBundleInstance = setupLocale()
	}
	return localeBundleInstance
}

func InitLocalizerInstance(locale language.Tag, accept language.Tag) *i18n.Localizer {
	localizerInstance = i18n.NewLocalizer(GetLocaleBundleInstance(), locale.String(), accept.String())
	return GetLocalizerInstance()
}

func GetLocalizerInstance() *i18n.Localizer {
	if localizerInstance == nil {
		panic("Singleton is not initialized. Call InitLocalizerInstance first.")
	}
	return localizerInstance
}

func setupLocale() *i18n.Bundle {
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
	_, err = bundle.LoadMessageFile("resources/en/validation.en.toml")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = bundle.LoadMessageFile("resources/ru/validation.ru.toml")
	if err != nil {
		fmt.Println(err.Error())
	}
	return bundle
}
