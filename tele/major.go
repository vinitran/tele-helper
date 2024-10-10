package tele

import (
	"fmt"
	"net/url"
	"strings"
)

type MajorApp struct {
	MiniappCfg
}

func NewMajorApp() *MajorApp {
	return &MajorApp{MiniappCfg{
		Url:        "https://web.telegram.org/k/#@major",
		UrlQueryId: "https://major.bot/",
		Name:       string(TonmarketAppName),
	}}
}

func (app *MajorApp) GetQueryId(input string) (string, error) {
	fragmentParts := strings.Split(input, "&")
	var tgWebAppData string
	for _, part := range fragmentParts {
		if strings.HasPrefix(part, "#tgWebAppData=") {
			tgWebAppData = strings.TrimPrefix(part, "#tgWebAppData=")
			break
		}
	}

	if tgWebAppData == "" {
		return "", fmt.Errorf("no tgWebAppData found")
	}

	decodedData, err := url.QueryUnescape(tgWebAppData)
	if err != nil {
		return "", fmt.Errorf("error decoding tgWebAppData: %v", err)
	}

	return decodedData, nil
}

func (app *MajorApp) GetUrl() string {
	return app.Url
}

func (app *MajorApp) GetUrlQueryId() string {
	return app.UrlQueryId
}

func (app *MajorApp) NameApp() string {
	return app.Name
}
