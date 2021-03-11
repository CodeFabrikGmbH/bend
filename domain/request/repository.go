package request

type Repository interface {
	Close()
	Save(req Request) error
	GetPaths() []string
	GetRequestCountForPath(path string) int
	GetRequestsForPath(path string) []Request
	GetRequest(path string, id string) Request

	DeletePath(path string) error
	DeleteRequestForPath(path string, id string) error
}
