package translate

import (
	"github.com/leonelquinteros/gotext"
	_ "golang.org/x/text/message/catalog"
)

//var local *gotext.Locale

func Config(lan_local string) {

	// Create Locale with library path and language code
	gotext.Configure("locals", lan_local, "messages")
//
	gotext.SetLanguage(lan_local)
//	po := new(gotext.Po)
//	// Parse .po file
//	po.ParseFile("locals/fa_IR/LC_MESSAGES/messages.po")
//
//	str := `
//msgid "One apple"
//msgstr "Una manzana"`
//
//	po.Parse([]byte(str))
//
//	// Get translation
//	println(po.Get("One apple"))
//	println(gotext.Get("Email"))
}

func setup(locale string, domain string, dir string) {

}
