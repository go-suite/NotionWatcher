package event

import "fmt"

// Type of event
type Type string

// Enum of types of events
const (
	PageAddedToDatabase   Type = "pageAddedToDatabase"
	PageUpdatedInDatabase Type = "pageUpdatedInDatabase"
)

func (c Type) String() string {
	return string(c)
}

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

// Event sent to the webhook
type Event struct {
	Name     Type     `json:"event"`
	Database Database `json:"database"`
	Page     Page     `json:"page"`
}

type Database struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Page struct {
	Id string `json:"id"`
}
