package test

import (
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
