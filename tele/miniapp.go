package tele

import (
	"errors"
)

type AppName string

var BlumAppName AppName = "blum"
var MajorAppName AppName = "major"

type Miniapp interface {
	GetQueryId(input string) (string, error)
	GetUrl() string
	GetUrlQueryId() string
	NameApp() string
}

type MiniappCfg struct {
	Url        string
	UrlQueryId string
}

func NewMiniapp(app string) (Miniapp, error) {
	switch AppName(app) {
	case BlumAppName:
		return NewBlumApp(), nil
	case MajorAppName:
		return NewMajorApp(), nil
	default:
		return nil, errors.New("app not supported")
	}
}
