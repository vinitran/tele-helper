package test

import (
	"fmt"
	"strings"
	"testing"

	"go-login/utils/file"
)

func TestReadFile(t *testing.T) {
	path := "./data_example.txt"
	data, err := file.ReadLines(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal("err: can not read file")
	}
}

func TestFolderExist(t *testing.T) {
	path := "./example"
	isExist := file.FolderExists(path)
	if !isExist {
		t.Fatal("err: can not find folder")
	}
}

func TestQueryId(t *testing.T) {
	// Test case
	testInput := "#tgWebAppData=query_id%3DAAHliQpOAgAAAOWJCk5h_C3f%26user%3D%257B%2522id%2522%253A5604280805%252C%2522first_name%2522%253A%2522Vini%2522%252C%2522last_name%2522%253A%2522Tran%2522%252C%2522username%2522%253A%2522Vinitrannn%2522%252C%2522language_code%2522%253A%2522en%2522%252C%2522allows_write_to_pm%2522%253Atrue%257D%26auth_date%3D1728536006%26hash%3D1b4e5b869c9c326b4e9566ed848256814e178c1248fbca43b34ab68f6e0ae146&tgWebAppVersion=7.10&tgWebAppPlatform=web"

	// Call the function
	queryID, err := ExtractQueryID(testInput)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Extracted query_id:", queryID)
	}
}

func ExtractQueryID(input string) (string, error) {
	queryIDStart := strings.Index(input, "query_id")
	if queryIDStart == -1 {
		return "", fmt.Errorf("query_id not found")
	}

	// Find the next '&' character after the query_id to determine where it ends
	queryIDEnd := strings.Index(input[queryIDStart:], "&")
	if queryIDEnd == -1 {
		// If there's no '&' after query_id, take the rest of the string
		return input[queryIDStart:], nil
	}

	// Extract the query_id segment
	return input[queryIDStart : queryIDStart+queryIDEnd], nil
}
