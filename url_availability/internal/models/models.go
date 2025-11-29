package models

type LinkStatus string

const (
	Available    LinkStatus = "available"
	NotAvailable LinkStatus = "not available"
)

type LinkTask struct {
	TaskNum  int                   `json:"task_num"`
	Links    []string              `json:"links"`
	Statuses map[string]LinkStatus `json:"statuses"`
	Done     bool                  `json:"done"`
}
