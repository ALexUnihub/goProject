package journal

import "net/http"

type JournalItem struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}

type Graph struct {
	Nodes         map[string][]string `json:"nodes"`
	ActivityIndex map[string]string   `json:"activity"`
}

//go:generate mockgen -source=journal.go -destination=repo_mock.go -package=journal JournalRepo
type JournalRepo interface {
	AddItemInRepo(r *http.Request) (*JournalItem, error)
	CreateGrapgh() (*Graph, error)
}
