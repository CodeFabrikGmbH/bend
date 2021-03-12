package application

import (
	"code-fabrik.com/bend/domain/dashboard"
	"code-fabrik.com/bend/domain/environment"
	"sort"
)

type DashboardService struct {
	Env environment.Environment
}

func (ds DashboardService) GetDashboardData(path string, requestId *string) dashboard.DashBoard {
	requests := ds.getRequests(path)
	if requestId == nil && len(requests) > 0 {
		requestId = &requests[0].ID
	}

	return dashboard.DashBoard{
		Paths:          ds.getPaths(),
		CurrentPath:    path,
		Requests:       requests,
		RequestDetails: ds.getRequestDetails(path, requestId),
	}
}

func (ds DashboardService) getPaths() []dashboard.Path {
	var requestPaths []dashboard.Path

	paths := ds.Env.RequestRepository.GetPaths()
	for _, p := range paths {
		count := ds.Env.RequestRepository.GetRequestCountForPath(p)

		requestPaths = append(requestPaths, dashboard.Path{
			Path:  p,
			Count: count,
		})
	}

	sort.SliceStable(requestPaths, func(i, j int) bool {
		return requestPaths[i].Count > requestPaths[j].Count
	})

	return requestPaths
}

func (ds DashboardService) getRequests(path string) []dashboard.Request {
	requests := ds.Env.RequestRepository.GetRequestsForPath(path)
	sort.SliceStable(requests, func(i, j int) bool {
		return requests[i].Timestamp > requests[j].Timestamp
	})

	var dashboardRequests []dashboard.Request
	for _, r := range requests {
		dashboardRequests = append(dashboardRequests, dashboard.CreateRequest(r))
	}
	return dashboardRequests
}

func (ds DashboardService) getRequestDetails(path string, id *string) dashboard.RequestDetails {
	if id == nil {
		return dashboard.RequestDetails{}
	}

	req := ds.Env.RequestRepository.GetRequest(path, *id)
	return dashboard.CreateRequestDetails(req)
}
