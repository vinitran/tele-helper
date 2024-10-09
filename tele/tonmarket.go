package tele

import (
	"fmt"
	"net/url"
	"strings"
)

type TonmarketApp struct {
	MiniappCfg
}

func NewTonmarketApp() *BlumApp {
	return &BlumApp{MiniappCfg{
		Url:        "https://web.telegram.org/k/#@Tomarket_ai_bot",
		UrlQueryId: "https://mini-app.tomarket.ai/",
		Name:       string(TonmarketAppName),
	}}
}

func (app *TonmarketApp) GetQueryId(input string) (string, error) {
	// Step 1: Extract the part containing tgWebAppData
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

func (app *TonmarketApp) GetUrl() string {
	return app.Url
}

func (app *TonmarketApp) GetUrlQueryId() string {
	return app.UrlQueryId
}

func (app *TonmarketApp) NameApp() string {
	return app.Name
}
