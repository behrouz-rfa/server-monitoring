package protos

type Protos interface {
	// Use to filter by keyword
	GetSearchKey() string
}

type Result struct {
	Completed bool
	Ok        bool
}
