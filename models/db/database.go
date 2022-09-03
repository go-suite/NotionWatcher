package db

// Database : Represent a notion database
// ref: https://developers.notion.com/reference/database
type Database struct {
	ID   uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID string `json:"uuid" gorm:"not null;unique"`
	Name string `json:"name" gorm:"not null"`
}
