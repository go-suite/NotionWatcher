package common

import (
	"context"
	"github.com/dstotijn/go-notion"
	wdb "github.com/gennesseaux/NotionWatcher/common/db"
	"github.com/gennesseaux/NotionWatcher/common/event"
	"github.com/gennesseaux/NotionWatcher/common/webhook"
	nwDatasource "github.com/gennesseaux/NotionWatcher/setup/datasource"
	log "github.com/go-mods/zerolog-quick"
	"gorm.io/gorm/clause"
	"time"
)

// datasource : datasource
var datasource = nwDatasource.Datasource

func (w *Watcher) sendWebHook(database *wdb.Database, id string) (err error) {
	log.Debug().Msgf("Sending webhook for %s to url: %s", w.Name, w.WebHook)

	evt := event.Event{
		Name: w.Type,
		Database: event.Database{
			Id:   database.UUID,
			Name: database.Title,
		},
		Page: event.Page{
			Id: id,
		},
	}

	return webhook.SendMessage(w.WebHook, evt)
}

func (w *Watcher) prepareDatabaseWatcher() (err error) {

	// Database object
	database := &wdb.Database{UUID: w.DatabaseId}
	datasource.DB().Find(&database, database)

	// Database watcher objects
	databaseWatcher := &wdb.DatabaseWatcher{Name: w.Name}
	datasource.DB().Preload(clause.Associations).Find(&databaseWatcher, databaseWatcher)

	// Update the database object
	database.UUID = w.DatabaseId
	database.Title = w.DatabaseName
	err = datasource.DB().Save(&database).Error
	if err != nil {
		return err
	}

	// Update the database watcher object
	databaseWatcher.Database = *database
	databaseWatcher.DatabaseId = database.ID
	err = datasource.DB().Save(&databaseWatcher).Error
	if err != nil {
		return err
	}

	return nil
}

func (w *Watcher) runDatabaseWatcher() (err error) {

	// Get the latest information stored in DatabaseWatcher table
	dw := &wdb.DatabaseWatcher{Name: w.Name}
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
	var hasMore = true

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

	// Get the latest record
	client := notion.NewClient(w.Token)
	resp, err := client.QueryDatabase(context.Background(), dw.Database.UUID, query)
	if err != nil {
		return err
	}

	//  There might be no records
	if len(resp.Results) == 0 {
		return nil
	}

	// Test if something changed since the last check
	if len(resp.Results) > 0 && (dw.LastRecordProcessed != resp.Results[0].ID || !dw.LastTimeChecked.Time.After(resp.Results[0].LastEditedTime)) {
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
	dw.LastRecordProcessed = resp.Results[0].ID
	err = datasource.DB().Save(&dw).Error
	if err != nil {
		return err
	}

	return nil
}
