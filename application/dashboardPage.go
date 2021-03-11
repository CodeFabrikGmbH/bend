package application

import (
	"code-fabrik.com/bend/domain/dashboardpage"
	"code-fabrik.com/bend/domain/environment"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type DashboardPage struct {
	Env environment.Environment
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

func (rs DashboardPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
		}
	}()

	_, err := rs.Env.Authentication.Authenticate(w, r)
	if err != nil {
		http.Redirect(w, r, "/login?origin="+r.RequestURI, http.StatusTemporaryRedirect)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/dashboard")
	requests := rs.getRequests(path)

	requestId := getQueryValueOrNil(r.URL.Query(), "requestId")
	if requestId == nil && len(requests) > 0 {
		requestId = &requests[0].ID
	}

	dashBoardData := dashboardpage.DashBoard{
		Paths:          rs.getPaths(),
		CurrentPath:    path,
		Requests:       requests,
		RequestDetails: rs.getRequestDetails(path, requestId),
	}

	rs.Env.DashboardPage.Present(w, dashBoardData)
}

func (rs DashboardPage) getPaths() []dashboardpage.Path {
	var requestPaths []dashboardpage.Path

	paths := rs.Env.RequestRepository.GetPaths()
	for _, p := range paths {
		count := rs.Env.RequestRepository.GetRequestCountForPath(p)

		requestPaths = append(requestPaths, dashboardpage.Path{
			Path:  p,
			Count: count,
		})
	}

	sort.SliceStable(requestPaths, func(i, j int) bool {
		return requestPaths[i].Count > requestPaths[j].Count
	})

	return requestPaths
}

func (rs DashboardPage) getRequests(path string) []dashboardpage.Request {
	requests := rs.Env.RequestRepository.GetRequestsForPath(path)
	sort.SliceStable(requests, func(i, j int) bool {
		return requests[i].Timestamp > requests[j].Timestamp
	})

	var dashboardRequests []dashboardpage.Request
	for _, r := range requests {
		dashboardRequests = append(dashboardRequests, dashboardpage.CreateRequest(r))
	}
	return dashboardRequests
}

func (rs DashboardPage) getRequestDetails(path string, id *string) dashboardpage.RequestDetails {
	if id == nil {
		return dashboardpage.RequestDetails{}
	}

	req := rs.Env.RequestRepository.GetRequest(path, *id)
	return dashboardpage.CreateRequestDetails(req)
}
