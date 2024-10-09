package tele

import (
	"fmt"
	"net/url"
	"strings"
)

type BlumApp struct {
	MiniappCfg
}

func NewBlumApp() *BlumApp {
	return &BlumApp{MiniappCfg{
		Url:        "https://web.telegram.org/k/#@BlumCryptoBot",
		UrlQueryId: "https://telegram.blum.codes/",
		Name:       string(BlumAppName),
	}}
}

func (app *BlumApp) GetQueryId(input string) (string, error) {
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

func (app *BlumApp) GetUrl() string {
	return app.Url
}

func (app *BlumApp) GetUrlQueryId() string {
	return app.UrlQueryId
}

func (app *BlumApp) NameApp() string {
	return app.Name
}
