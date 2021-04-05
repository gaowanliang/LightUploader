package main

import (
	"io/ioutil"
	"log"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.SimplifiedChinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	rd, err := ioutil.ReadDir("i18n")
	if err != nil {
		log.Panic(err)
	}
	for _, fi := range rd {
		if !fi.IsDir() && path.Ext(fi.Name()) == ".toml" {
			bundle.LoadMessageFile("i18n/" + fi.Name())
		}
	}
	/*bundle.MustParseMessageFileBytes([]byte(

			), "zh-CN.toml")
		bundle.MustParseMessageFileBytes([]byte(`
	HelloWorld = "Hello World!"
	`), "en.toml")*/

}

type Loc struct {
	localize *i18n.Localizer
}

func (loc *Loc) init(tag string) {
	loc.localize = i18n.NewLocalizer(bundle, tag)
}
func (loc *Loc) print(tag string) string {
	return loc.localize.MustLocalize(&i18n.LocalizeConfig{MessageID: tag})
}
