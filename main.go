package main

import (
	"fmt"
	"log"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	templateFileName := config.Excel.TemplateFileName

	sheetCoordinates := ExcelCoordinates{
		sheetName:       config.Excel.SheetName,
		startingRow:     config.Excel.SheetCoordinates.StartingRow,
		cellHours:       config.Excel.SheetCoordinates.CellHours,
		cellLocations:   config.Excel.SheetCoordinates.CellLocations,
		cellDescription: config.Excel.SheetCoordinates.CellDescription,
		cellComment:     config.Excel.SheetCoordinates.CellComment,
	}

	header := Header{
		firstname:   config.Excel.Firstname,
		surname:     config.Excel.Surname,
		company:     config.Excel.Company,
		team:        config.Excel.Team,
		month:       config.Excel.Month,
		year:        config.Excel.Year,
		sheetName:   config.Excel.SheetName,
		startingRow: config.Excel.HeaderCoordinates.StartingRow,
		column:      config.Excel.HeaderCoordinates.Column,
	}

	harvestConfig := HarvestApiConfig{
		accountId: config.Harvest.AccountId,
		apiKey:    config.Harvest.ApiKey,
		url:       config.Harvest.Url,
	}

	newFileName := fmt.Sprintf("TS_EDP_%d_%02d_%s_%s_%s.xlsx", header.year, header.month, header.company, header.surname, header.firstname)

	f, err := openFile(templateFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	timeEntries, err := getAllTimeEntries(harvestConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	projectIds := HarvestProjectIds{
		projectId50hz:       config.Harvest.ProjectId50hz,
		projectIdSickleave:  config.Harvest.ProjectIdSickleave,
		projectIdPubHoliday: config.Harvest.ProjectIdPubHoliday,
	}

	filteredTimeEntries := filterTimeEntries(timeEntries, projectIds)
	mappedTimeEntries := mapTimeEntriesToDay(filteredTimeEntries)
	tasks := createTaskFromTimeEntries(mappedTimeEntries, config.Excel.Month)

	for _, task := range tasks {
		addTask(f, task, sheetCoordinates)
	}

	addHeader(f, header)

	if err := saveFile(f, newFileName); err != nil {
		fmt.Println(err)
	}

	reEvaluateFormulae(newFileName)

	fmt.Println("File saved to:", newFileName)
}
