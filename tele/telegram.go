package tele

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/xuri/excelize/v2"
	log2 "go-login/utils/log"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"go-login/utils/file"
)

const (
	SRC_USER_DIR = "./config/chrome"
	TELEGRAM_URL = "https://web.telegram.org/"
)

type Client struct {
	Config
	Miniapp
	fileWrite *excelize.File
	log       *log2.LogHelper
}

func NewClient(app *string, cfg Config, fileWrite *excelize.File) (*Client, error) {
	if !cfg.isValid() {
		return nil, errors.New("invalid config")
	}

	logHelper := log2.NewLogHelper(cfg.Name, cfg.Proxy, "")
	if app == nil {
		return &Client{
			cfg,
			nil,
			nil,
			logHelper,
		}, nil
	}
	miniApp, err := NewMiniapp(*app)
	if err != nil {
		return nil, logHelper.ErrorMessage(err)
	}

	logHelper = log2.NewLogHelper(cfg.Name, cfg.Proxy, miniApp.NameApp())
	return &Client{cfg, miniApp, fileWrite, logHelper}, nil
}

func (c *Client) Login() error {
	c.log.Success("start login")
	extensionPath := "config/extensions/gleekbfjekiniecknbkamfmkohkpodhe"

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("load-extension", extensionPath),
		chromedp.UserDataDir(SRC_USER_DIR),
		chromedp.Flag("profile-directory", c.Name),
		// chromedp.Headless,
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-logging", true),
		chromedp.WindowSize(560, 1080),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	}

	//if c.useProxy() {
	//	opts = append(opts, chromedp.ProxyServer(c.Proxy))
	//}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(ctx, network.Enable())
	if err != nil {
		return c.log.ErrorMessage(err)
	}
	log.Println("asdasd")
	err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(TELEGRAM_URL),
		chromedp.WaitVisible(`//*[@id="LeftColumn-main"]`, chromedp.BySearch),
		chromedp.Sleep(3 * time.Second),
	})
	if err != nil {
		return c.log.ErrorMessage(err)
	}

	c.log.Success("login successfully")
	return nil
}

func (c *Client) GetDataTele() (string, error) {
	c.log.Success("start get query_id")

	extensionPath := "config/extensions/gleekbfjekiniecknbkamfmkohkpodhe"
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("load-extension", extensionPath),
		chromedp.UserDataDir(SRC_USER_DIR),
		chromedp.Flag("profile-directory", c.Name),
		// chromedp.Headless,
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-logging", true),
		chromedp.WindowSize(480, 1080),
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoFirstRun,
	}

	if c.useProxy() {
		opts = append(opts, chromedp.ProxyServer(c.Proxy))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(ctx, network.Enable())
	if err != nil {
		return "", c.log.ErrorMessage(err)
	}

	var telegramData string
	chromedp.Run(ctx,
		network.Enable(),
	)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if eventRequest, ok := ev.(*network.EventRequestWillBeSent); ok {
			if eventRequest.Request.URL == c.GetUrlQueryId() {
				telegramData, err = c.GetQueryId(eventRequest.Request.URLFragment)
				if err != nil {
					log.Println(c.log.ErrorMessage(err))
				}
			}
		}
	})

	startBtn := `/html/body/div[1]/div/div[2]/div/div[1]/div[4]/div/div[1]/div/div[8]/div[1]/div[2]`
	runBtn := "/html/body/div[6]/div/div[2]/button[1]"

	err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(c.GetUrl()),
		chromedp.Sleep(7 * time.Second),
		chromedp.WaitVisible(startBtn, chromedp.BySearch),
		chromedp.Click(startBtn, chromedp.BySearch),
		chromedp.Sleep(5 * time.Second),

		// runBtn
		chromedp.ActionFunc(func(ctx context.Context) error {
			timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
			defer cancel()

			if err := chromedp.Run(timeoutCtx, chromedp.WaitVisible(runBtn, chromedp.BySearch)); err != nil {
				return err
			}

			err = chromedp.Run(
				timeoutCtx,
				chromedp.Sleep(2*time.Second),
				chromedp.Click(runBtn, chromedp.BySearch),
				chromedp.Sleep(12*time.Second),
			)
			if err != nil {
				return err
			}
			return nil
		}),
	})
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return "", c.log.ErrorMessage(err)
	}

	if telegramData == "" {
		return "", c.log.ErrorWithMsg("can not get query_id")
	}

	c.log.Success("get query_id successfully")
	return telegramData, nil
}

func (c *Client) ExportQueryId(row int) error {
	telegramData, err := c.GetDataTele()
	if err != nil {
		return err
	}

	err = file.ExportDataToExel(c.fileWrite, row, c.Name, c.Proxy, telegramData)
	if err != nil {
		return c.log.ErrorMessage(err)
	}

	c.log.Success("export query_id to excel file successfully")
	return nil
}
