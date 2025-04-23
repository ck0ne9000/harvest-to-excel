package main

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type Header struct {
	firstname   string
	surname     string
	company     string
	team        string
	month       int
	year        int
	sheetName   string
	startingRow int
	column      string
}

type ExcelCoordinates struct {
	sheetName       string
	startingRow     int
	cellHours       string
	cellLocations   string
	cellDescription string
	cellComment     string
}

type Task struct {
	day         int
	hours       float64
	location    string
	description string
	comment     string
}

func closeFile(f *excelize.File) error {
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func openFile(fileName string) (*excelize.File, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return f, err
	}
	defer closeFile(f)
	return f, nil
}

func addTask(f *excelize.File, task Task, excelCoordinates ExcelCoordinates) {
	activeRow := strconv.Itoa(excelCoordinates.startingRow + task.day - 1)

	f.SetCellValue(excelCoordinates.sheetName, fmt.Sprintf("%s%s", excelCoordinates.cellHours, activeRow), task.hours)
	f.SetCellValue(excelCoordinates.sheetName, fmt.Sprintf("%s%s", excelCoordinates.cellLocations, activeRow), task.location)
	f.SetCellValue(excelCoordinates.sheetName, fmt.Sprintf("%s%s", excelCoordinates.cellDescription, activeRow), task.description)
	f.SetCellValue(excelCoordinates.sheetName, fmt.Sprintf("%s%s", excelCoordinates.cellComment, activeRow), task.comment)
}

func addHeader(f *excelize.File, header Header) {
	f.SetCellValue(header.sheetName, fmt.Sprintf("%s%d", header.column, header.startingRow), header.firstname+header.surname)
	f.SetCellValue(header.sheetName, fmt.Sprintf("%s%d", header.column, header.startingRow+1), header.company)
	f.SetCellValue(header.sheetName, fmt.Sprintf("%s%d", header.column, header.startingRow+2), header.team)

	fmt.Println("Please rember to change the month in the dropdown, due to macros this can't be done with this tool :(")
}

func reEvaluateFormulae(fileName string) {
	f, err := openFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	f.UpdateLinkedValue()

	if err := saveFile(f, fileName); err != nil {
		fmt.Println(err)
	}
}

func saveFile(f *excelize.File, fileName string) error {
	if err := f.SaveAs(fileName); err != nil {
		return err
	}
	return nil
}
