package models

// Database : Represent a db database
// ref: https://developers.notion.com/reference/database
//
// UUID: database id
// Title: database title
type Database struct {
	ID    uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID  string `json:"uuid" gorm:"not null;unique"`
	Title string `json:"name" gorm:"not null"`
}
