package main

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Harvest struct {
		AccountId           int    `mapstructure:"accountId"`
		ApiKey              string `mapstructure:"apiKey"`
		ProjectId50hz       int    `mapstructure:"projectId50hz"`
		ProjectIdSickleave  int    `mapstructure:"projectIdSickleave"`
		ProjectIdPubHoliday int    `mapstructure:"projectIdPubHoliday"`
		Url                 string `mapstructure:"url"`
	} `mapstructure:"harvest"`
	Excel struct {
		TemplateFileName  string `mapstructure:"templateFileName"`
		Company           string `mapstructure:"company"`
		Firstname         string `mapstructure:"firstname"`
		Surname           string `mapstructure:"surname"`
		Team              string `mapstructure:"team"`
		Month             int    `mapstructure:"month"`
		Year              int    `mapstructure:"year"`
		SheetName         string `mapstructure:"sheetName"`
		HeaderCoordinates struct {
			StartingRow int    `mapstructure:"startingRow"`
			Column      string `mapstructure:"column"`
		} `mapstructure:"HeaderCoordinates"`
		SheetCoordinates struct {
			StartingRow     int    `mapstructure:"startingRow"`
			CellHours       string `mapstructure:"cellHours"`
			CellLocations   string `mapstructure:"cellLocations"`
			CellDescription string `mapstructure:"cellDescription"`
			CellComment     string `mapstructure:"cellComment"`
		} `mapstructure:"SheetCoordinates"`
	} `mapstructure:"excel"`
}

func loadConfig() (Config, error) {
	viper.SetConfigName("harvest_exporter")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$PWD")
	viper.AddConfigPath("$XDG_CONFIG_HOME/harvest_exporter")

	viper.SetDefault("Harvest.ProjectIdSickleave", 42796995)
	viper.SetDefault("Harvest.ProjectIdPubHoliday", 42765148)
	viper.SetDefault("Excel.SheetName", "Time Sheet EDP")
	viper.SetDefault("Excel.TemplateFileName", "TS_EDP_yyyy_mm_company_surname_1stname1.xlsx")
	viper.SetDefault("Excel.HeaderCoordinates.StartingRow", 2)
	viper.SetDefault("Excel.HeaderCoordinates.Column", "E")
	viper.SetDefault("Excel.SheetCoordinates.StartingRow", 8)
	viper.SetDefault("Excel.SheetCoordinates.CellHours", "C")
	viper.SetDefault("Excel.SheetCoordinates.CellLocations", "D")
	viper.SetDefault("Excel.SheetCoordinates.CellDescription", "E")
	viper.SetDefault("Excel.SheetCoordinates.cellComment", "F")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	year := time.Now().Year()
	month := time.Now().Month()

	pflag.IntVar(&config.Excel.Month, "month", int(month), "month to create the timesheet for")
	pflag.IntVar(&config.Excel.Year, "year", year, "year to create the timesheet for")
	pflag.Parse()

	return config, nil
}
