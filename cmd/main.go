package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"go-login/tele"
	"go-login/utils/file"
	"log"
	"sync"
)

const maxConcurrency = 4

func main() {
	users, err := file.ReadFileExcel("./data/input.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	////login
	//for _, user := range users[1:] {
	//	teleCli, err := tele.NewClient(
	//		tele.BlumAppName,
	//		tele.Config{
	//			Name:  user[0],
	//			Proxy: user[1],
	//		},
	//	)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	err = teleCli.Login()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}

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
				log.Fatal(err)
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
}
