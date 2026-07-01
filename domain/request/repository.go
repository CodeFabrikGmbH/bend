package request

type Repository interface {
	Add(req Request) (Request, error)
	GetPaths() []string
	GetPathCounts() map[string]int
	GetRequestCountForPath(path string) int
	GetRequestsForPath(path string) []Request
	GetSummariesForPath(path string) []Summary
	GetSummariesPage(path string, before int, limit int) ([]Summary, bool)
	GetRequest(path string, id string) Request

	DeletePath(path string) error
	DeleteRequestForPath(path string, id string) error
}
