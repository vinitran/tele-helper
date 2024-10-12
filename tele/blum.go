package tele

type BlumApp struct {
	MiniappCfg
}

func NewBlumApp() *BlumApp {
	return &BlumApp{MiniappCfg{
		Url:  "https://web.telegram.org/k/#?tgaddr=tg%3A%2F%2Fresolve%3Fdomain%3Dblum%26appname%3Dapp%26startapp%3Dref_ILdTUAZ37i",
		Name: string(BlumAppName),
	}}
}

func (app *BlumApp) GetUrl() string {
	return app.Url
}

func (app *BlumApp) NameApp() string {
	return app.Name
}
