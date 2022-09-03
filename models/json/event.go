package json

// Event used in webhook

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
