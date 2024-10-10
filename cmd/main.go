package main

import (
	"fmt"
	"go-login/exporter"
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

			err = teleCli.ExportQueryId(i + 1)
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
	var userProfiles []string

	arr.ArrEach(users, func(user []string) {
		userProfiles = append(userProfiles, fmt.Sprintf("./config/chrome/%s", user[0]))
	})

	backupFolder := "./data/usersData"

	// Step 1: Copy all user profiles to the backup folder
	err = exporter.SaveAllUserProfilesToFolder(userProfiles, backupFolder)
	if err != nil {
		return fmt.Errorf("failed to copy user profiles: %v", err)
	}

	// Step 2: Zip the entire backup folder
	destinationZip := "./data/usersDataExport.zip"
	err = exporter.ZipFolder(backupFolder, destinationZip)
	if err != nil {
		return fmt.Errorf("failed to zip the backup folder: %v", err)
	}

	log.Println("All user profiles successfully saved and zipped.")
	return nil
}

func importUserData(c *cli.Context) error {
	zipFilePath := "./data/"              // Path to your zip file
	chromeProfileDir := "./config/chrome" // Path to Chrome's profile directory

	// Import user profiles from the zip file
	err := exporter.ImportUserProfilesFromZip(zipFilePath, chromeProfileDir)
	if err != nil {
		return fmt.Errorf("Failed to import user profiles: %v", err)
	}

	return nil
}
