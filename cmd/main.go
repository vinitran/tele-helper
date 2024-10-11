package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"go-login/utils/arr"

	"github.com/urfave/cli/v2"
	"github.com/xuri/excelize/v2"
	"go-login/tele"
	"go-login/utils/file"
)

const maxConcurrency = 2

var (
	NameFlag = &cli.StringFlag{
		Name:  "name",
		Usage: "The name of the person to greet",
	}
	MaxConcurencyFlag = &cli.StringFlag{
		Name:  "threads",
		Usage: "The number of concurrent threads to use",
	}
	AppFlag = &cli.StringFlag{
		Name:     "app",
		Usage:    "App name: blum major",
		Required: true,
	}
)

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
			Flags:   append(flags, AppFlag, NameFlag, MaxConcurencyFlag),
		},
		{
			Name:    "export",
			Aliases: []string{},
			Usage:   "export to zip user data",
			Action:  exportUserData,
			Flags:   append(flags, NameFlag),
		},
		{
			Name:    "import",
			Aliases: []string{},
			Usage:   "import from zip user data",
			Action:  importUserData,
			Flags:   append(flags),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func getUsers(c *cli.Context) ([][]string, error) {
	userInput, err := file.ReadFileExcel("./data/input.xlsx")
	if err != nil {
		return nil, err
	}
	// remove first row
	userInput = userInput[1:]

	specificName := c.String("name")
	if specificName == "" {
		return userInput, nil
	}

	var users [][]string
	usernames := strings.Split(specificName, ",")
	arr.ArrEach(userInput, func(user []string) {
		arr.ArrEach(usernames, func(s string) {
			if user[0] == s {
				users = append(users, user)
			}
		})
	})
	return users, nil
}

func login(c *cli.Context) error {
	users, err := getUsers(c)
	if err != nil {
		return err
	}

	for _, user := range users {
		teleCli, err := tele.NewClient(
			nil,
			tele.Config{
				Name:  user[0],
				Proxy: user[1],
			},
			nil,
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
	users, err := getUsers(c)
	if err != nil {
		return err
	}

	sem := make(chan struct{}, maxConcurrency)
	if c.Int("threads") != 0 {
		sem = make(chan struct{}, c.Int("threads"))
	}

	app := c.String("app")

	var wg sync.WaitGroup
	fileWrite := excelize.NewFile()

	for index, user := range users {
		if len(user) < 2 {
			continue
		}
		wg.Add(1)
		go func(i int, nm, prx string, wg *sync.WaitGroup, sem chan struct{}) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			teleCli, err := tele.NewClient(
				&app,
				tele.Config{
					Name:  nm,
					Proxy: prx,
				},
				fileWrite,
			)
			if err != nil {
				log.Println(err)
				return
			}

			err = teleCli.ExportQueryId(context.Background(), i+1)
			if err != nil {
				log.Println(err)
				return
			}
		}(index, user[0], user[1], &wg, sem)
	}
	wg.Wait()

	fmt.Println("All windows opened with max concurrency control.")
	return nil
}

func exportUserData(c *cli.Context) error {
	users, err := getUsers(c)
	if err != nil {
		return err
	}

	backupFolderRoot := "./data/usersData"

	arr.ArrEachWithErr(users, func(user []string) error {
		userProfile := fmt.Sprintf("./config/profiles/%s/Default", user[0])
		backupFolder := fmt.Sprintf("%s/%s", backupFolderRoot, user[0])
		err := file.CopyFolder(userProfile, backupFolder)
		if err != nil {
			return err
		}
		return nil
	})

	destinationZip := "./data/profiles.zip"

	err = file.ZipFolder(backupFolderRoot, destinationZip)
	if err != nil {
		return err
	}

	log.Println("All user profiles successfully saved and zipped.")

	return file.DeleteFolder(backupFolderRoot)
}

func importUserData(c *cli.Context) error {
	zipFilePath := "./data/" // Path to your zip file
	zipExtract := fmt.Sprintf("%s%s", zipFilePath, "zip_extract")
	chromeProfileDir := "./config/profiles" // Path to Chrome's profile directory

	err := file.UnzipAllFilesInFolder(zipFilePath, zipExtract)
	if err != nil {
		return err
	}

	users, err := file.GetFoldersInFolder(zipExtract)
	if err != nil {
		return err
	}

	err = arr.ArrEachWithErr(users, func(user string) error {
		userProfileImport := fmt.Sprintf("%s/%s", zipExtract, user)
		err := file.CopyFolder(
			fmt.Sprintf("%s/example", chromeProfileDir),
			fmt.Sprintf("%s/%s", chromeProfileDir, user),
		)
		if err != nil {
			return err
		}

		userPath := fmt.Sprintf("%s/%s/Default", chromeProfileDir, user)
		err = file.CopyFolder(userProfileImport, userPath)
		if err != nil {
			return err
		}
		log.Printf("imported user: %s\n", user)

		return nil
	})
	if err != nil {
		return err
	}

	return file.DeleteFolder(zipExtract)

}
