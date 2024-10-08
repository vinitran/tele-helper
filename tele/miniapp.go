package tele

import (
	"errors"
)

type AppName string

var BlumAppName AppName = "BLUM"

type Miniapp interface {
	GetQueryId(input string) (string, error)
	GetUrl() string
	GetUrlQueryId() string
}

type MiniappCfg struct {
	Url        string
	UrlQueryId string
}

func NewMiniapp(app AppName) (Miniapp, error) {
	switch app {
	case BlumAppName:
		return NewBlumApp(), nil
	default:
		return nil, errors.New("app not supported")
	}
}
