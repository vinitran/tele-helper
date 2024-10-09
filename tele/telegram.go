package tele

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"go-login/utils/file"
)

const (
	SRC_USER_DIR = "./config/example"
	TELEGRAM_URL = "https://web.telegram.org/"
)

type Client struct {
	Config
	Miniapp
}

func NewClient(app AppName, cfg Config) (*Client, error) {
	if !cfg.isValid() {
		return nil, errors.New("invalid config")
	}

	miniApp, err := NewMiniapp(app)
	if err != nil {
		return nil, err
	}
	return &Client{cfg, miniApp}, nil
}

func (c *Client) Login() error {
	userDir := fmt.Sprintf("./config/%s", c.Name)
	err := file.CheckExistAndCopy(userDir, SRC_USER_DIR)
	if err != nil {
		return err
	}

	opts := []chromedp.ExecAllocatorOption{
		chromedp.UserDataDir(userDir),
		// chromedp.Headless,
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-logging", true),
		chromedp.WindowSize(480, 1080),
		chromedp.ProxyServer(c.Proxy),
	}

	if c.useProxy() {
		opts = append(opts, chromedp.ProxyServer(c.Proxy))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err = chromedp.Run(ctx, network.Enable())
	if err != nil {
		return err
	}

	err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(TELEGRAM_URL),
		chromedp.Sleep(5 * time.Second),
		chromedp.WaitVisible(`//*[@id="page-chats"]`, chromedp.BySearch),
		chromedp.Sleep(50 * time.Second),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetDataTele() (string, error) {
	userDir := fmt.Sprintf("./config/%s", c.Name)
	err := file.CheckExistAndCopy(userDir, SRC_USER_DIR)
	if err != nil {
		return "", err
	}

	opts := []chromedp.ExecAllocatorOption{
		chromedp.UserDataDir(userDir),
		// chromedp.Headless,
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-logging", true),
		chromedp.WindowSize(480, 1080),
		chromedp.ProxyServer(c.Proxy),
	}

	if c.useProxy() {
		opts = append(opts, chromedp.ProxyServer(c.Proxy))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err = chromedp.Run(ctx, network.Enable())
	if err != nil {
		return "", err
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
					log.Println(err)
				}
			}
		}
	})

	startBtn := `/html/body/div[1]/div/div[2]/div/div[1]/div[4]/div/div[1]/div/div[8]/div[1]/div[2]`
	runBtn := "/html/body/div[7]/div/div[2]/button[1]"

	err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(c.GetUrl()),
		chromedp.Sleep(10 * time.Second),
		chromedp.WaitVisible(startBtn, chromedp.BySearch),
		chromedp.Click(startBtn, chromedp.BySearch),
		chromedp.Sleep(10 * time.Second),

		// runBtn
		chromedp.ActionFunc(func(ctx context.Context) error {
			timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := chromedp.Run(timeoutCtx, chromedp.WaitVisible(runBtn, chromedp.BySearch)); err != nil {
				fmt.Println("runBtn not found within 5 seconds, proceeding...")
				return nil
			}

			return chromedp.Click(runBtn, chromedp.BySearch).Do(ctx)
		}),
	})
	if err != nil {
		return "", err
	}

	return telegramData, nil
}
