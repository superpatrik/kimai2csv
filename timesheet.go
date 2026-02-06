package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"timesheet/kimai"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type appConfig struct {
	apiToken   string
	url        string
	projectIds []string
	user       string
	begin      string
	end        string
	extended   bool
	csv        bool
}

type DateEntriesContainer struct {
	Date    kimai.KimaiTime
	Entries map[string]ProjectEntriesContainer
}

type ProjectEntriesContainer struct {
	Project kimai.Project
	Entries map[string]ActivityEntriesContainer
}

type ActivityEntriesContainer struct {
	Activity kimai.Activity
	Entries  []kimai.TimeSheet
}

func main() {
	apiTokenParam := flag.String("apiToken", "", "Kimai api token")
	urlParam := flag.String("url", "", "Kimai url")
	projectsParam := flag.String("projects", "", "Kimai projects")
	userParam := flag.String("user", "", "Kimai user")
	beginParam := flag.String("begin", "", "Begin of timespan")
	endParam := flag.String("end", "", "End of timespan")
	extendedParam := flag.Bool("extended", false, "Extended output, shows activities")
	csvParam := flag.String("csv", "", "Output file for csv")
	lastMonthParam := flag.Bool("lastMonth", false, "Use last month for begin and end")

	flag.Parse()

	if *apiTokenParam == "" {
		log.Fatal("Kimai api token required")
	}

	if *urlParam == "" {
		log.Fatal("Kimai url required")
	}

	if *beginParam == "" && *endParam == "" {
		now := time.Now()

		if *lastMonthParam {
			startOfLastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
			startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			endOfLastMonth := startOfMonth.Add(-time.Nanosecond)

			*beginParam = startOfLastMonth.Format("2006-01-02T15:04:05")
			*endParam = endOfLastMonth.Format("2006-01-02T15:04:05")
		} else {
			startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			startOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
			endOfMonth := startOfNextMonth.Add(-time.Nanosecond)

			*beginParam = startOfMonth.Format("2006-01-02T15:04:05")
			*endParam = endOfMonth.Format("2006-01-02T15:04:05")
		}
	}

	var projectIds []string

	if *projectsParam != "" {
		projectIds = strings.Split(*projectsParam, ",")
	}

	config := appConfig{
		apiToken:   *apiTokenParam,
		url:        *urlParam,
		projectIds: projectIds,
		user:       *userParam,
		begin:      *beginParam,
		end:        *endParam,
	}

	client := kimai.NewClient(config.url, config.apiToken)

	projects := client.Projects()
	activities := client.Activities()
	timeSheets := client.TimeSheets(config.user, config.begin, config.end, config.projectIds, 500)

	projectMap := make(map[string]kimai.Project)
	activityMap := make(map[string]kimai.Activity)

	for _, project := range projects {
		projectId := strconv.Itoa(project.Id)
		projectMap[projectId] = project
	}

	for _, activity := range activities {
		activityId := strconv.Itoa(activity.Id)
		activityMap[activityId] = activity
	}

	entries := make(map[string]DateEntriesContainer)

	var dates []string

	for _, timeSheet := range timeSheets {
		date := timeSheet.Begin.Format("2006-01-02")
		projectId := strconv.Itoa(timeSheet.ProjectId)
		activityId := strconv.Itoa(timeSheet.ActivityId)

		dateEntries, ok := entries[date]

		if !ok {
			dateEntries = DateEntriesContainer{
				Date:    timeSheet.Begin,
				Entries: make(map[string]ProjectEntriesContainer),
			}
		}

		projectEntries, ok := dateEntries.Entries[projectId]

		if !ok {
			projectEntries = ProjectEntriesContainer{
				Project: projectMap[projectId],
				Entries: make(map[string]ActivityEntriesContainer),
			}
		}

		activityEntries, ok := projectEntries.Entries[activityId]

		if !ok {
			activityEntries = ActivityEntriesContainer{
				Activity: activityMap[activityId],
				Entries:  []kimai.TimeSheet{},
			}
		}

		activityEntries.Entries = append(activityEntries.Entries, timeSheet)
		projectEntries.Entries[activityId] = activityEntries
		dateEntries.Entries[projectId] = projectEntries
		entries[date] = dateEntries

		if !slices.Contains(dates, date) {
			dates = append(dates, date)
		}
	}

	slices.Sort(dates)

	hours := 0.0
	amount := 0.0

	for _, date := range dates {
		formattedDate := entries[date].Date.Format("02.01.2006")

		fmt.Printf("🗓️ %s\n", formattedDate)

		entriesContainer := entries[date]

		for _, projectEntriesContainer := range entriesContainer.Entries {
			fmt.Printf("\t 📌 %s\n", projectEntriesContainer.Project.Name)

			for _, activityEntriesContainer := range projectEntriesContainer.Entries {
				if *extendedParam {
					fmt.Printf("\t\t 📋 %s\n", activityEntriesContainer.Activity.Name)
				}

				for _, timeSheet := range activityEntriesContainer.Entries {
					begin := timeSheet.Begin.Format("15:04")
					end := timeSheet.End.Format("15:04")

					if *extendedParam {
						fmt.Printf("\t\t\t 🕖 %s -> %s: %s\n", begin, end, timeSheet.Description)
					} else {
						fmt.Printf("\t\t 🕖 %s -> %s: %s\n", begin, end, timeSheet.Description)
					}

					hours += float64(timeSheet.Duration) / 3600.0
					amount += timeSheet.Rate

				}
			}
		}
	}

	start, _ := time.Parse("2006-01-02T15:04:05", *beginParam)
	end, _ := time.Parse("2006-01-02T15:04:05", *endParam)

	fmt.Printf("\n")
	fmt.Printf("Start: %s\n", start.Format("02.01.2006, 15:04"))
	fmt.Printf("Ende: %s\n", end.Format("02.01.2006, 15:04"))
	fmt.Printf("\n")

	p := message.NewPrinter(language.German)
	p.Printf("🕖 %.2f Stunden -> %.2f € / %.2f (brutto)\n", hours, amount, amount*1.19)

	if *csvParam != "" {
		f, err := os.Create(*csvParam)

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		if *extendedParam {
			f.WriteString("Datum;Beginn;Ende;Dauer;Projekt;Tätigkeit;Beschreibung;Preis\n")
		} else {
			f.WriteString("Datum;Beginn;Ende;Dauer;Projekt;Beschreibung;Preis\n")
		}

		for _, date := range dates {
			formattedDate := entries[date].Date.Format("02.01.2006")
			entriesContainer := entries[date]

			for _, projectEntriesContainer := range entriesContainer.Entries {
				if *extendedParam {
					for _, activityEntriesContainer := range projectEntriesContainer.Entries {
						for _, timeSheet := range activityEntriesContainer.Entries {
							begin := timeSheet.Begin.Format("15:04")
							end := timeSheet.End.Format("15:04")

							activityId := strconv.Itoa(timeSheet.ActivityId)
							duration := float64(timeSheet.Duration) / 3600

							line := fmt.Sprintf(
								"%s;%s;%s;%.2f;%s;%s;%s;%.2f €\n",
								formattedDate,
								begin,
								end,
								duration,
								projectEntriesContainer.Project.Name,
								activityMap[activityId].Name,
								timeSheet.Description,
								timeSheet.Rate,
							)

							f.WriteString(line)
						}
					}
				} else {
					doings := make([]string, 0)
					duration := 0.0
					groupedAmount := 0.0
					now := time.Now()
					begin := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())

					for _, activityEntriesContainer := range projectEntriesContainer.Entries {
						for _, timeSheet := range activityEntriesContainer.Entries {
							parts := strings.Split(timeSheet.Description, ",")

							for _, part := range parts {
								doing := strings.TrimSpace(part)

								if !slices.Contains(doings, doing) {
									doings = append(doings, doing)
								}
							}

							duration += float64(timeSheet.Duration)
							groupedAmount += timeSheet.Rate
						}
					}

					end := begin.Add(time.Duration(duration) * time.Second)

					line := fmt.Sprintf(
						"%s;%s;%s;%.2f;%s;%s;%.2f €\n",
						formattedDate,
						begin.Format("15:04"),
						end.Format("15:04"),
						duration/3600,
						projectEntriesContainer.Project.Name,
						strings.Join(doings, ", "),
						groupedAmount,
					)

					f.WriteString(line)
				}
			}
		}
	}
}
