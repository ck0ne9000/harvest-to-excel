package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HarvestApiConfig struct {
	url       string
	accountId int
	apiKey    string
}

type HarvestProjectIds struct {
	projectId50hz       int
	projectIdSickleave  int
	projectIdPubHoliday int
}

type TimeEntries struct {
	ID                int64       `json:"id"`
	SpentDate         string      `json:"spent_date"`
	Hours             float64     `json:"hours"`
	HoursWithoutTimer float64     `json:"hours_without_timer"`
	RoundedHours      float64     `json:"rounded_hours"`
	Notes             string      `json:"notes"`
	IsLocked          bool        `json:"is_locked"`
	LockedReason      interface{} `json:"locked_reason"`
	IsClosed          bool        `json:"is_closed"`
	IsBilled          bool        `json:"is_billed"`
	TimerStartedAt    interface{} `json:"timer_started_at"`
	StartedTime       interface{} `json:"started_time"`
	EndedTime         interface{} `json:"ended_time"`
	IsRunning         bool        `json:"is_running"`
	Billable          bool        `json:"billable"`
	Budgeted          bool        `json:"budgeted"`
	BillableRate      interface{} `json:"billable_rate"`
	CostRate          interface{} `json:"cost_rate"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	User              struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Client struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Currency string `json:"currency"`
	} `json:"client"`
	Project struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"project"`
	Task struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"task"`
	UserAssignment struct {
		ID               int         `json:"id"`
		IsProjectManager bool        `json:"is_project_manager"`
		IsActive         bool        `json:"is_active"`
		UseDefaultRates  bool        `json:"use_default_rates"`
		Budget           interface{} `json:"budget"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		HourlyRate       interface{} `json:"hourly_rate"`
	} `json:"user_assignment"`
	TaskAssignment struct {
		ID         int         `json:"id"`
		Billable   bool        `json:"billable"`
		IsActive   bool        `json:"is_active"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
		HourlyRate interface{} `json:"hourly_rate"`
		Budget     float64     `json:"budget"`
	} `json:"task_assignment"`
	Invoice           interface{} `json:"invoice"`
	ExternalReference interface{} `json:"external_reference"`
}

type JsonResponse struct {
	TimeEntries  []TimeEntries `json:"time_entries"`
	PerPage      int           `json:"per_page"`
	TotalPages   int           `json:"total_pages"`
	TotalEntries int           `json:"total_entries"`
	NextPage     interface{}   `json:"next_page"`
	PreviousPage interface{}   `json:"previous_page"`
	Page         int           `json:"page"`
	Links        struct {
		First    string      `json:"first"`
		Next     interface{} `json:"next"`
		Previous interface{} `json:"previous"`
		Last     string      `json:"last"`
	} `json:"links"`
}

type groupedTimeEntries map[string][]TimeEntries

var jsonResponse JsonResponse

func getAllTimeEntries(config HarvestApiConfig) (JsonResponse, error) {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, config.url, nil)
	if err != nil {
		return JsonResponse{}, err
	}
	req.Header.Add("Harvest-Account-Id", strconv.Itoa(config.accountId))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.apiKey))

	res, err := client.Do(req)
	if err != nil {
		return JsonResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return JsonResponse{}, err
	}

	json.Unmarshal(body, &jsonResponse)

	return jsonResponse, nil
}

func filterTimeEntries(jsonResponse JsonResponse, projectIds HarvestProjectIds) []TimeEntries {
	var filteredTimeEntries []TimeEntries

	for _, v := range jsonResponse.TimeEntries {
		if v.Project.ID == projectIds.projectId50hz {
			filteredTimeEntries = append(filteredTimeEntries, v)
		}

		if v.Project.ID == projectIds.projectIdSickleave || v.Project.ID == projectIds.projectIdPubHoliday {
			v.Hours = 0.0
			filteredTimeEntries = append(filteredTimeEntries, v)
		}
	}
	return filteredTimeEntries
}

func mapTimeEntriesToDay(timeEntries []TimeEntries) groupedTimeEntries {
	mappedTimeEntries := make(map[string][]TimeEntries)

	for _, v := range timeEntries {
		mappedTimeEntries[v.SpentDate] = append(mappedTimeEntries[v.SpentDate], v)
	}
	return mappedTimeEntries
}

func createTaskFromTimeEntries(timeEntries groupedTimeEntries, billingMonth int) []Task {
	var tasks []Task

	for day, entries := range timeEntries {
		hours := 0.0
		description := ""
		comment := ""
		location := ""

		dateElements := strings.Split(day, "-")
		monthString := dateElements[1]

		month, err := strconv.Atoi(monthString)
		if err != nil {
			fmt.Println("Conversion error:", err)
		}

		if month != billingMonth {
			continue
		}

		dayString := dateElements[2]

		day, err := strconv.Atoi(dayString)
		if err != nil {
			fmt.Println("Conversion error:", err)
		}

		for _, entry := range entries {
			hours += entry.Hours
			noteElements := strings.Split(entry.Notes, ": ")
			description = description + noteElements[0] + "\n"

			if len(noteElements) == 1 {
				continue
			}

			comment = comment + noteElements[1] + "\n"
		}

		if hours != 0.0 {
			location = "remote"
		}

		task := Task{
			day:         day,
			hours:       hours,
			location:    location,
			description: description,
			comment:     comment,
		}
		tasks = append(tasks, task)
	}
	return tasks
}
