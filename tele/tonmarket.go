package tele

type TonmarketApp struct {
	MiniappCfg
}

func NewTonmarketApp() *TonmarketApp {
	return &TonmarketApp{MiniappCfg{
		Url:  "https://web.telegram.org/k/#?tgaddr=tg%3A%2F%2Fresolve%3Fdomain%3DTomarket_ai_bot%26appname%3Dapp%26startapp%3D0002tXzw",
		Name: string(TonmarketAppName),
	}}
}

func (app *TonmarketApp) GetUrl() string {
	return app.Url
}

func (app *TonmarketApp) NameApp() string {
	return app.Name
}
