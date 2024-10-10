package tele

import (
	"fmt"
	"strings"
)

type TonmarketApp struct {
	MiniappCfg
}

func NewTonmarketApp() *TonmarketApp {
	return &TonmarketApp{MiniappCfg{
		Url:        "https://web.telegram.org/k/#@Tomarket_ai_bot",
		UrlQueryId: "https://mini-app.tomarket.ai/",
		Name:       string(TonmarketAppName),
	}}
}

func (app *TonmarketApp) GetQueryId(input string) (string, error) {
	queryIDStart := strings.Index(input, "query_id")
	if queryIDStart == -1 {
		return "", fmt.Errorf("query_id not found")
	}

	// Find the next '&' character after the query_id to determine where it ends
	queryIDEnd := strings.Index(input[queryIDStart:], "&")
	if queryIDEnd == -1 {
		// If there's no '&' after query_id, take the rest of the string
		return input[queryIDStart:], nil
	}

	// Extract the query_id segment
	return input[queryIDStart : queryIDStart+queryIDEnd], nil
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
