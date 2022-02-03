package utils

import (
	"github.com/gorilla/sessions"
	"net"
	"server-monitoring/shared/consts"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math/rand"
	"strings"
	"time"
)

var (
	Helper helpersInterface = &helper{}
)

type helpersInterface interface {
	RandStringRunes(n int) string
}
type helper struct {
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (h helper) RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return strings.ToLower(string(b))
}

func PriceFormat(price float32) string {
	p := message.NewPrinter(language.Persian)
	return p.Sprintf("%d", int(price))
}

func LanguageCurrency(sess *sessions.Session) (int, int) {
	languageId := 1
	currency := 1
	if sess.Values[consts.LANGUAGE] != nil {
		languageId = sess.Values[consts.LANGUAGE].(int)
	}
	if sess.Values[consts.CURRENCY] != nil {
		currency = sess.Values[consts.CURRENCY].(int)
	}
	return languageId, currency
}

func IsIpV4(ip string) bool {
	ipConv := net.ParseIP(ip)
	if len(ipConv) == net.IPv4len {
		return true
	}
	return false
}
