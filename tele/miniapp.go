package tele

import (
	"errors"
)

type AppName string

var (
	BlumAppName      AppName = "blum"
	MajorAppName     AppName = "major"
	TonmarketAppName AppName = "tonmarket"
)

type Miniapp interface {
	GetUrl() string
	NameApp() string
}

type MiniappCfg struct {
	Url  string
	Name string
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
