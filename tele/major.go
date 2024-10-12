package tele

type MajorApp struct {
	MiniappCfg
}

func NewMajorApp() *MajorApp {
	return &MajorApp{MiniappCfg{
		Url:  "https://web.telegram.org/k/#?tgaddr=tg%3A%2F%2Fresolve%3Fdomain%3Dmajor%26appname%3Dstart",
		Name: string(TonmarketAppName),
	}}
}

func (app *MajorApp) GetUrl() string {
	return app.Url
}

func (app *MajorApp) NameApp() string {
	return app.Name
}
