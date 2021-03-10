package application

import (
	"code-fabrik.com/bend/domain/dashboard"
	"code-fabrik.com/bend/domain/request"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type DashboardService struct {
	RequestRepository  request.Repository
	DashboardPresenter dashboard.Presenter
}

func getQueryValueOrNil(v url.Values, key string) *string {
	value := v.Get(key)
	if len(value) == 0 {
		value = v.Get(strings.ToLower(key))
		if len(value) == 0 {
			return nil
		}
	}
	return &value
}

func (rs DashboardService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
		}
	}()
	path := strings.TrimPrefix(r.URL.Path, "/dashboard")
	requests := rs.getRequests(path)

	requestId := getQueryValueOrNil(r.URL.Query(), "requestId")
	if requestId == nil && len(requests) > 0 {
		requestId = &requests[0].ID
	}

	dashBoardData := dashboard.DashBoard{
		Paths:          rs.getPaths(),
		CurrentPath:    path,
		Requests:       requests,
		RequestDetails: rs.getRequestDetails(path, requestId),
	}

	rs.DashboardPresenter.Present(w, dashBoardData)
}

func (rs DashboardService) getPaths() []dashboard.Path {
	var requestPaths []dashboard.Path

	paths := rs.RequestRepository.GetPaths()
	for _, p := range paths {
		count := rs.RequestRepository.GetRequestCountForPath(p)

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

func (rs DashboardService) getRequests(path string) []dashboard.Request {
	requests := rs.RequestRepository.GetRequestsForPath(path)
	sort.SliceStable(requests, func(i, j int) bool {
		return requests[i].Timestamp > requests[j].Timestamp
	})

	var dashboardRequests []dashboard.Request
	for _, r := range requests {
		dashboardRequests = append(dashboardRequests, dashboard.CreateRequest(r))
	}
	return dashboardRequests
}

func (rs DashboardService) getRequestDetails(path string, id *string) dashboard.RequestDetails {
	if id == nil {
		return dashboard.RequestDetails{}
	}

	req := rs.RequestRepository.GetRequest(path, *id)
	return dashboard.CreateRequestDetails(req)
}
