package models

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

// DatabaseWatcher : keep tracks of the latest checked records
type DatabaseWatcher struct {
	ID                  uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                string         `json:"name" gorm:"not null;unique"`
	DatabaseId          uint64         `json:"database_id"`
	Database            Database       `json:"database"`
	LastTimeChecked     sql.NullTime   `json:"last_time_checked" gorm:"type:TIMESTAMP NULL"`
	LastRecordProcessed string         `json:"last_uuid_checked"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at"`
}

func (w *DatabaseWatcher) GetLastTimeChecked(now time.Time) time.Time {
	if w.LastTimeChecked.Valid {
		return w.LastTimeChecked.Time
	} else {
		return now
	}
}
