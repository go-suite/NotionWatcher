package json

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	dbModels "github.com/gennesseaux/NotionWatcher/models/db"
	nwConfig "github.com/gennesseaux/NotionWatcher/setup/config"
	nwDatasource "github.com/gennesseaux/NotionWatcher/setup/datasource"
	nwLogger "github.com/gennesseaux/NotionWatcher/setup/logger"

	"github.com/dstotijn/go-notion"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

// logger : logger
var logger = nwLogger.Logger

// config : config
var config = nwConfig.Config

// datasource : datasource
var datasource = nwDatasource.Datasource

// Type used in watcher
type Type string

// Enum of types
const (
	PageAddedToDatabase   Type = "pageAddedToDatabase"
	PageUpdatedInDatabase Type = "pageUpdatedInDatabase"
)

// Watcher : watcher json file
type Watcher struct {
	Name       string     `json:"name" validate:"required"`
	Type       Type       `json:"type" validate:"required,oneof=pageAddedToDatabase pageUpdatedInDatabase"`
	DatabaseId string     `json:"database_id" validate:"required"`
	StartDate  *time.Time `json:"start_date"`
	WebHook    string     `json:"webhook" validate:"required,url"`
	Cron       string     `json:"cron" validate:"required"`
}

func (w Watcher) GetSortTimestampType() notion.SortTimestamp {
	if w.Type == PageAddedToDatabase {
		return notion.SortTimeStampCreatedTime
	} else if w.Type == PageUpdatedInDatabase {
		return notion.SortTimeStampLastEditedTime
	}
	return ""
}

func (w Watcher) GetStartDate() *time.Time {
	if w.StartDate != nil {
		startDate := *w.StartDate
		startDate = startDate.Truncate(time.Minute)
		return &startDate
	}
	return nil
}

func (w Watcher) PageIsBefore(p notion.Page, d time.Time) bool {
	if w.Type == PageAddedToDatabase {
		return p.CreatedTime.Before(d)
	} else if w.Type == PageUpdatedInDatabase {
		return p.LastEditedTime.Before(d)
	}
	return true
}

func (w Watcher) PageIsSameOrBefore(p notion.Page, d time.Time) bool {
	if w.Type == PageAddedToDatabase {
		return p.CreatedTime.Equal(d) || p.CreatedTime.Before(d)
	} else if w.Type == PageUpdatedInDatabase {
		return p.LastEditedTime.Equal(d) || p.LastEditedTime.Before(d)
	}
	return true
}

func (w Watcher) sendWebHook(database *dbModels.Database, id string) (err error) {
	logger.Debug(fmt.Sprintf("Sending webhook for %s", w.Name))
	logger.Debug(fmt.Sprintf("%s", w.WebHook))
	event := Event{
		Name: w.Type,
		Database: Database{
			Id:   database.UUID,
			Name: database.Name,
		},
		Page: Page{
			Id: id,
		},
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("%s", jsonData))

	request, err := http.NewRequest("GET", w.WebHook, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Error sending hook : %s", resp.Status))
	}
	return nil
}

func (w Watcher) install() (err error) {
	// Database
	if w.Type == "pageAddedToDatabase" || w.Type == "pageUpdatedInDatabase" {
		err = w.installDatabaseWatcher()
		if err != nil {
			return err
		}
	}

	// Pages

	// Blocks

	return nil
}

func (w Watcher) installDatabaseWatcher() (err error) {

	// Get notion database
	client := notion.NewClient(config.Token)
	ndb, err := client.FindDatabaseByID(context.Background(), w.DatabaseId)
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("Found notion database with id : %s", ndb.ID))
	logger.Debug(fmt.Sprintf("  --> database name : %s", ndb.Title[0].PlainText))

	// Database object
	database := &dbModels.Database{UUID: ndb.ID}
	datasource.DB().Find(&database, database)

	// Database watcher objects
	databaseWatcher := &dbModels.DatabaseWatcher{Name: w.Name}
	datasource.DB().Preload(clause.Associations).Find(&databaseWatcher, databaseWatcher)

	// Update the database object
	database.UUID = ndb.ID
	database.Name = ndb.Title[0].PlainText
	err = datasource.DB().Save(&database).Error
	if err != nil {
		logger.Error("Error saving to database", zap.Error(err))
	}

	// Update the database watcher object
	databaseWatcher.Database = *database
	databaseWatcher.DatabaseId = database.ID
	err = datasource.DB().Save(&databaseWatcher).Error
	if err != nil {
		logger.Error("Error saving to database", zap.Error(err))
	}

	return nil
}

func (w Watcher) Run() (err error) {
	logger.Debug(fmt.Sprintf("Running watcher : %s", w.Name))

	// Install
	err = w.install()
	if err != nil {
		return err
	}

	// Database
	if w.Type == "pageAddedToDatabase" || w.Type == "pageUpdatedInDatabase" {
		err = w.RunDatabaseWatcher()
		if err != nil {
			logger.Error(fmt.Sprintf("Error running watcher : %s", w.Name), zap.Error(err))
		}
	}

	// Pages

	// Blocks

	return nil
}

func (w Watcher) RunDatabaseWatcher() (err error) {

	// Get the latest information stored in DatabaseWatcher table
	dw := &dbModels.DatabaseWatcher{Name: w.Name}
	err = datasource.DB().Preload(clause.Associations).Find(&dw, dw).Error
	if err != nil {
		return err
	}

	// dates
	now := time.Now().Truncate(time.Minute) // <-- Notion don't use seconds
	startDate := dw.GetLastTimeChecked(now).Truncate(time.Minute)
	endDate := now

	// the start date also depends on the date of the watcher
	if w.GetStartDate() != nil {
		if w.GetStartDate().After(startDate) {
			startDate = *w.GetStartDate()
		}
	}

	// Array to store all the records to process
	var results []notion.Page

	//
	var hasMore bool = true

	// Query sent to notion
	query := &notion.DatabaseQuery{
		PageSize: 1,
		Sorts: []notion.DatabaseQuerySort{
			{
				Timestamp: w.GetSortTimestampType(),
				Direction: notion.SortDirDesc,
			},
		},
	}

	// Get last record
	client := notion.NewClient(config.Token)
	resp, err := client.QueryDatabase(context.Background(), dw.Database.UUID, query)
	if err != nil {
		return err
	}

	//  There might be no records
	if len(resp.Results) == 0 {
		return nil
	}

	// Test if something changed since the last check
	if len(resp.Results) > 0 && (dw.LastRecordProccesed != resp.Results[0].ID || !dw.LastTimeChecked.Time.After(resp.Results[0].LastEditedTime)) {
		for {
			query.PageSize = 60
			resp, err = client.QueryDatabase(context.Background(), dw.Database.UUID, query)
			if err != nil {
				return err
			}

			for _, result := range resp.Results {
				if !w.PageIsBefore(result, startDate) {
					results = append(results, result)
				}
			}

			if len(results) == 0 {
				break
			}

			hasMore = resp.HasMore

			if resp.NextCursor != nil {
				query.StartCursor = *resp.NextCursor
			}

			// While condition
			lastResult := results[len(results)-1]
			cont := !w.PageIsSameOrBefore(lastResult, startDate) && hasMore
			if !(cont) {
				break
			}
		}
	}

	// Process results
	for _, result := range results {
		// Call the web hook
		err = w.sendWebHook(&dw.Database, result.ID)
		if err != nil {
			return err
		}
	}

	// Save last processed
	dw.LastTimeChecked.Time = endDate
	dw.LastTimeChecked.Valid = true
	dw.LastRecordProccesed = resp.Results[0].ID
	err = datasource.DB().Save(&dw).Error
	if err != nil {
		return err
	}

	return nil
}

func (w Watcher) Watch(s *gocron.Scheduler) (err error) {
	_, err = s.Cron(w.Cron).Do(w.Run)
	if err != nil {
		return err
	}
	s.StartAsync()
	return
}
