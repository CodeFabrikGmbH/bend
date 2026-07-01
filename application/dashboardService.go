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

type DashBoardViewData struct {
	Paths          []Path
	CurrentPath    string
	Requests       []RequestAbstract
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
	requests := ds.getRequests(path)
	if len(requestId) == 0 && len(requests) > 0 {
		requestId = requests[0].ID
	}

	return DashBoardViewData{
		Paths:          ds.getPaths(),
		CurrentPath:    path,
		Requests:       requests,
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

func (ds DashboardService) getRequests(path string) []RequestAbstract {
	summaries := ds.RequestRepository.GetSummariesForPath(path)
	sort.SliceStable(summaries, func(i, j int) bool {
		return summaries[i].Timestamp > summaries[j].Timestamp
	})

	var dashboardRequests []RequestAbstract
	for _, s := range summaries {
		dashboardRequests = append(dashboardRequests, RequestAbstract{
			ID:        s.ID,
			Timestamp: time.Unix(0, s.Timestamp).Format("2 Jan 2006 15:04:05"),
			Method:    s.Method,
			Status:    s.Response.ResponseStatusCode,
		})
	}
	return dashboardRequests
}

func (ds DashboardService) getRequestDetails(path string, id string) RequestDetails {
	if len(id) == 0 {
		return RequestDetails{}
	}

	req := ds.RequestRepository.GetRequest(path, id)
	return CreateRequestDetails(req)
}
