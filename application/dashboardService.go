package application

import (
	"code-fabrik.com/bend/domain/request"
	"sort"
	"time"
)

type Path struct {
	Path  string
	Count int
}

type RequestAbstract struct {
	ID        string `json:"ID"`
	Timestamp string `json:"timestamp"`
	Unix      int64  `json:"unix"` // milliseconds, for client-side relative times
	Method    string `json:"method"`
	Status    int    `json:"status"`
}

type RequestDetails struct {
	ID        string              `json:"ID"`
	Timestamp string              `json:"timestamp"`
	Path      string              `json:"path"`
	Method    string              `json:"method"`
	Body      string              `json:"body"`
	Header    map[string][]string `json:"header"`
	Host      string              `json:"host"`
	Uri       string              `json:"uri"`
	Response  request.Response    `json:"response"`
}

// dashboardPageSize is how many requests are rendered for a path on the first
// load. Older requests are fetched incrementally as the user scrolls, so a busy
// path (hundreds of thousands of requests) never has to be read into memory or
// serialized into the page at once.
const dashboardPageSize = 50

// dashboardPathLimit caps how many endpoints are rendered in the sidebar. With
// tens of thousands of distinct paths the full list would produce a multi-MB
// page; the most active paths are shown.
const dashboardPathLimit = 500

type DashBoardViewData struct {
	Paths          []Path
	PathTotal      int
	CurrentPath    string
	Requests       []RequestAbstract
	RequestTotal   int
	HasMore        bool
	RequestDetails RequestDetails
}

func CreateRequestDetails(request request.Request) RequestDetails {
	return RequestDetails{
		ID:        request.ID,
		Timestamp: time.Unix(0, request.Timestamp).Format("2 Jan 2006 15:04:05"),
		Path:      request.Path,
		Method:    request.Method,
		Body:      request.Body,
		Header:    request.Header,
		Host:      request.Host,
		Uri:       request.Uri,
		Response:  request.Response,
	}
}

type DashboardService struct {
	RequestRepository request.Repository
}

func (ds DashboardService) GenerateDashboardViewData(path string, requestId string) DashBoardViewData {
	requests, hasMore := ds.getRequests(path)
	total := ds.RequestRepository.GetRequestCountForPath(path)

	if len(requestId) == 0 && len(requests) > 0 {
		requestId = requests[0].ID
	}

	paths := ds.getPaths()
	pathTotal := len(paths)
	if len(paths) > dashboardPathLimit {
		paths = paths[:dashboardPathLimit]
	}

	return DashBoardViewData{
		Paths:          paths,
		PathTotal:      pathTotal,
		CurrentPath:    path,
		Requests:       requests,
		RequestTotal:   total,
		HasMore:        hasMore,
		RequestDetails: ds.getRequestDetails(path, requestId),
	}
}

func (ds DashboardService) getPaths() []Path {
	counts := ds.RequestRepository.GetPathCounts()

	requestPaths := make([]Path, 0, len(counts))
	for path, count := range counts {
		requestPaths = append(requestPaths, Path{
			Path:  path,
			Count: count,
		})
	}

	sort.SliceStable(requestPaths, func(i, j int) bool {
		if requestPaths[i].Count != requestPaths[j].Count {
			return requestPaths[i].Count > requestPaths[j].Count
		}
		return requestPaths[i].Path < requestPaths[j].Path
	})

	return requestPaths
}

func (ds DashboardService) getRequests(path string) ([]RequestAbstract, bool) {
	summaries, hasMore := ds.RequestRepository.GetSummariesPage(path, 0, dashboardPageSize)

	dashboardRequests := make([]RequestAbstract, 0, len(summaries))
	for _, s := range summaries {
		dashboardRequests = append(dashboardRequests, newRequestAbstract(s))
	}
	return dashboardRequests, hasMore
}

// newRequestAbstract projects a stored summary into the list-view shape used by
// both the initial page render and the incremental pagination API.
func newRequestAbstract(s request.Summary) RequestAbstract {
	return RequestAbstract{
		ID:        s.ID,
		Timestamp: time.Unix(0, s.Timestamp).Format("2 Jan 2006 15:04:05"),
		Unix:      s.Timestamp / int64(time.Millisecond),
		Method:    s.Method,
		Status:    s.Response.ResponseStatusCode,
	}
}

func (ds DashboardService) getRequestDetails(path string, id string) RequestDetails {
	if len(id) == 0 {
		return RequestDetails{}
	}

	req := ds.RequestRepository.GetRequest(path, id)
	return CreateRequestDetails(req)
}
