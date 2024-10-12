package tele

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	log2 "go-login/utils/log"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"go-login/utils/file"
)

const (
	SRC_USER_DIR = "./config/profiles/%s"
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

	opts := c.defaultOpt()
	if c.useProxy() {
		opts = append(opts, chromedp.ProxyServer(c.Proxy))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err := chromedp.Run(ctx, network.Enable())
	if err != nil {
		return c.log.ErrorMessage(err)
	}

	err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(TELEGRAM_URL),
		chromedp.WaitVisible(`//*[@id="LeftColumn-main"]`, chromedp.BySearch),
		chromedp.Sleep(20 * time.Second),
	})
	if err != nil {
		return c.log.ErrorMessage(err)
	}

	c.log.Success("login successfully")
	return nil
}

func (c *Client) GetDataTele(ctxCli context.Context) (string, error) {
	c.log.Success("start get query_id")

	opts := c.defaultOpt()
	if c.useProxy() {
		opts = append(opts, chromedp.ProxyServer(c.Proxy))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctxCli, opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	launchBtn := "/html/body/div[6]/div/div[2]/button[1]"
	var iframeSrc string

	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(c.GetUrl()),
		chromedp.WaitVisible(launchBtn, chromedp.BySearch),
		chromedp.Click(launchBtn, chromedp.BySearch),
		chromedp.WaitVisible(`iframe.payment-verification`, chromedp.ByQuery),
		chromedp.AttributeValue(`iframe.payment-verification`, "src", &iframeSrc, nil),
	})
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return "", c.log.ErrorMessage(err)
	}

	queryId, err := extractTgWebAppData(iframeSrc)
	if err != nil {
		return "", err
	}

	c.log.Success("get query_id successfully")
	return queryId, nil
}

func (c *Client) ExportQueryId(ctx context.Context, row int) error {
	telegramData, err := c.GetDataTele(ctx)
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

func extractTgWebAppData(src string) (string, error) {
	parts := strings.Split(src, "#")
	if len(parts) < 2 {
		return "", fmt.Errorf("no fragment found in the URL")
	}

	fragment := parts[1]
	fragmentParams := strings.Split(fragment, "&")

	for _, param := range fragmentParams {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) == 2 && kv[0] == "tgWebAppData" {
			return kv[1], nil
		}
	}

	return "", fmt.Errorf("tgWebAppData not found in URL fragment")
}

func (c *Client) defaultOpt() []chromedp.ExecAllocatorOption {
	extensionPath := "config/extensions/gleekbfjekiniecknbkamfmkohkpodhe"
	return []chromedp.ExecAllocatorOption{
		chromedp.Flag("load-extension", extensionPath),
		chromedp.UserDataDir(fmt.Sprintf(SRC_USER_DIR, c.Name)),
		// chromedp.Headless,
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-logging", true),
		chromedp.WindowSize(560, 1080),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	}
}
