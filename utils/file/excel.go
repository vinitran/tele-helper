package file

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadFileExcel(path string) ([][]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	return f.GetRows("Sheet1")
}

func ExportDataToExel(fileWrite *excelize.File, row int, name, proxy, queryId string) error {
	err := fileWrite.SetCellStr("Sheet1", fmt.Sprintf("A%d", row), name)
	if err != nil {
		return err
	}

	err = fileWrite.SetCellStr("Sheet1", fmt.Sprintf("B%d", row), proxy)
	if err != nil {
		return err
	}

	err = fileWrite.SetCellStr("Sheet1", fmt.Sprintf("C%d", row), queryId)
	if err != nil {
		return err
	}

	err = fileWrite.SaveAs("./data/output.xlsx")
	if err != nil {
		return err
	}
	return nil
}
