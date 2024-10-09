package tele

import (
	"errors"
)

type AppName string

var BlumAppName AppName = "blum"
var MajorAppName AppName = "major"
var TonmarketAppName AppName = "tonmarket"

type Miniapp interface {
	GetQueryId(input string) (string, error)
	GetUrl() string
	GetUrlQueryId() string
	NameApp() string
}

type MiniappCfg struct {
	Url        string
	UrlQueryId string
	Name       string
}

func NewMiniapp(app string) (Miniapp, error) {
	switch AppName(app) {
	case BlumAppName:
		return NewBlumApp(), nil
	case MajorAppName:
		return NewMajorApp(), nil
	case TonmarketAppName:
		return NewTonmarketApp(), nil
	default:
		return nil, errors.New("app not supported")
	}
}
