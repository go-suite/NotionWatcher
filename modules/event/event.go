package event

import "fmt"

// Event is the struct that represents an event
// It is used to send webhook events
type Event struct {
	Name     Type     `json:"event"`
	Database Database `json:"database"`
	Page     Page     `json:"page"`
}

// Database is the struct which contains the notion database information
type Database struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Page is the struct which contains the notion page information
type Page struct {
	Id string `json:"id"`
}

// Type of event
// It can be PageAddedToDatabase or PageUpdatedInDatabase
type Type string

// types of events
const (
	PageAddedToDatabase   Type = "pageAddedToDatabase"
	PageUpdatedInDatabase Type = "pageUpdatedInDatabase"
)

// ParseType parses a string into a Type
func ParseType(s string) (t Type, err error) {
	types := map[Type]struct{}{
		PageAddedToDatabase:   {},
		PageUpdatedInDatabase: {},
	}

	t = Type(s)
	_, ok := types[t]
	if !ok {
		return t, fmt.Errorf(`cannot parse:[%s] as type`, s)
	}
	return t, nil
}

func (c Type) String() string {
	return string(c)
}
