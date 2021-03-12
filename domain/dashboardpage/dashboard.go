package dashboardpage

import (
	"code-fabrik.com/bend/domain/request"
	"time"
)

type Path struct {
	Path  string
	Count int
}

type Request struct {
	ID        string `json:"ID"`
	Timestamp string `json:"timestamp"`
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

type DashBoard struct {
	Paths          []Path
	CurrentPath    string
	Requests       []Request
	RequestDetails RequestDetails
}

func CreateRequest(request request.Request) Request {
	return Request{
		ID:        request.ID,
		Timestamp: time.Unix(0, request.Timestamp).Format("2 Jan 2006 15:04:05"),
	}
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
