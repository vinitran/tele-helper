package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/xuri/excelize/v2"
	"go-login/tele"
	"go-login/utils/file"
)

const maxConcurrency = 2

var NameFlag = &cli.StringFlag{
	Name:  "name",
	Usage: "The name of the person to greet",
}

func main() {
	app := cli.NewApp()
	app.Name = "go-tele"
	flags := []cli.Flag{}
	app.Commands = []*cli.Command{
		{
			Name:    "login",
			Aliases: []string{},
			Usage:   "login",
			Action:  login,
			Flags:   append(flags, NameFlag),
		},
		{
			Name:    "queryid",
			Aliases: []string{},
			Usage:   "get query id",
			Action:  getQueryId,
			Flags:   append(flags),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func login(c *cli.Context) error {
	if c.String("name") != "" {
		teleCli, err := tele.NewClient(
			tele.BlumAppName,
			tele.Config{
				Name:  c.String("name"),
				Proxy: "",
			},
		)
		if err != nil {
			return err
		}

		err = teleCli.Login()
		if err != nil {
			return err
		}
		return nil
	}

	users, err := file.ReadFileExcel("./data/input.xlsx")
	if err != nil {
		return err
	}

	for _, user := range users[1:] {
		teleCli, err := tele.NewClient(
			tele.BlumAppName,
			tele.Config{
				Name:  user[0],
				Proxy: user[1],
			},
		)
		if err != nil {
			return err
		}

		err = teleCli.Login()
		if err != nil {
			return err
		}
	}
	return nil
}

func getQueryId(c *cli.Context) error {
	users, err := file.ReadFileExcel("./data/input.xlsx")
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)
	fileWrite := excelize.NewFile()

	for index, user := range users[1:] {
		if len(user) < 2 {
			continue
		}
		wg.Add(1)
		go func(i int, nm, prx string, wg *sync.WaitGroup, sem chan struct{}) {
			defer wg.Done()
			log.Println("tele account name", nm)
			// Acquire semaphore to limit concurrency
			sem <- struct{}{}
			defer func() { <-sem }() // Release semaphore when done

			teleCli, err := tele.NewClient(
				tele.BlumAppName,
				tele.Config{
					Name:  nm,
					Proxy: prx,
				},
			)
			if err != nil {
				return
			}

			telegramData, err := teleCli.GetDataTele()
			if err != nil {
				log.Fatal(err)
			}

			err = file.ExportDataToExel(fileWrite, i+1, nm, prx, telegramData)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Telegram data:", telegramData)
		}(index, user[0], user[1], &wg, sem)
	}
	wg.Wait()

	fmt.Println("All windows opened with max concurrency control.")
	return nil
}
