package watcher

import (
	"context"
	"errors"
	"github.com/dstotijn/go-notion"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/event"
	"github.com/go-playground/validator/v10"
	"time"
)

// config : instance of Config
var config = nwConfig.Config

// Watcher is a struct that represents a watcher
type Watcher struct {
	Name         string     `json:"name" validate:"required"`
	Type         event.Type `json:"type" validate:"required,oneof=pageAddedToDatabase pageUpdatedInDatabase"`
	DatabaseId   string     `json:"database_id" validate:"required"`
	DatabaseName string     `json:"database_name"`
	StartDate    *time.Time `json:"start_date"`
	WebHook      string     `json:"webhook"  validate:"required,url"`
	Cron         string     `json:"cron" validate:"omitempty"`
	Token        string     `json:"token" validate:"required"`
	Inactive     bool
}

// Validate validates the watcher,
// if the validation fails, an error is returned
func (w *Watcher) Validate() error {
	// Token
	if len(w.Token) == 0 {
		return errors.New("the token provided is not valid")
	}

	// DatabaseId
	if len(w.DatabaseId) == 0 {
		return errors.New("provide db database id")
	}

	// Get database
	if len(w.DatabaseId) > 0 {
		client := notion.NewClient(w.Token)
		ndb, err := client.FindDatabaseByID(context.Background(), w.DatabaseId)
		if err != nil {
			return err
		}
		w.DatabaseName = ndb.Title[0].PlainText
	}

	// Type
	if len(w.Type) == 0 {
		return errors.New("the type provided is not valid")
	}

	// Cron todo
	if len(w.Cron) == 0 {
		return errors.New("the cron provided is not valid")
	}

	// Validate
	validate := validator.New()
	return validate.Struct(w)
}

func (w *Watcher) GetSortTimestampType() notion.SortTimestamp {
	if w.Type == event.PageAddedToDatabase {
		return notion.SortTimeStampCreatedTime
	} else if w.Type == event.PageUpdatedInDatabase {
		return notion.SortTimeStampLastEditedTime
	}
	return ""
}

func (w *Watcher) GetStartDate() *time.Time {
	if w.StartDate != nil {
		startDate := *w.StartDate
		startDate = startDate.Truncate(time.Minute)
		return &startDate
	}
	return nil
}

func (w *Watcher) PageIsBefore(p notion.Page, d time.Time) bool {
	if w.Type == event.PageAddedToDatabase {
		return p.CreatedTime.Before(d)
	} else if w.Type == event.PageUpdatedInDatabase {
		return p.LastEditedTime.Before(d)
	}
	return true
}

func (w *Watcher) PageIsSameOrBefore(p notion.Page, d time.Time) bool {
	if w.Type == event.PageAddedToDatabase {
		return p.CreatedTime.Equal(d) || p.CreatedTime.Before(d)
	} else if w.Type == event.PageUpdatedInDatabase {
		return p.LastEditedTime.Equal(d) || p.LastEditedTime.Before(d)
	}
	return true
}
