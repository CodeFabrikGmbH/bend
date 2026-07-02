package application

import (
	"code-fabrik.com/bend/domain/config"
	"code-fabrik.com/bend/domain/request"
	"time"
)

type RequestService struct {
	RequestRepository request.Repository
	ConfigRepository  config.Repository
	Transport         request.Transport
	Hub               *EventHub
}

const defaultStatusCode = 200

// RequestPage is a page of request summaries returned by the pagination API.
type RequestPage struct {
	Items   []RequestAbstract `json:"items"`
	HasMore bool              `json:"hasMore"`
}

// ListRequests returns a page of request summaries for a path, newest first.
// before is an exclusive upper bound on the request id (0 for the newest page).
func (rs RequestService) ListRequests(path string, before int, limit int) RequestPage {
	summaries, hasMore := rs.RequestRepository.GetSummariesPage(path, before, limit)

	items := make([]RequestAbstract, 0, len(summaries))
	for _, s := range summaries {
		items = append(items, newRequestAbstract(s))
	}
	return RequestPage{Items: items, HasMore: hasMore}
}

func (rs RequestService) DeleteAllRequestsForPath(path string) error {
	return rs.RequestRepository.DeletePath(path)
}

func (rs RequestService) DeleteRequest(path string, requestId string) error {
	return rs.RequestRepository.DeleteRequestForPath(path, requestId)
}

func (rs RequestService) DeleteAllRequests() error {
	return rs.RequestRepository.DeleteAllRequests()
}

func (rs RequestService) SendRequestToTarget(path, requestId, targetUrl string) request.Response {
	req := rs.RequestRepository.GetRequest(path, requestId)
	return rs.Transport.SendRequestToTarget(req, targetUrl)
}

func (rs RequestService) TrackRequest(req request.Request) request.Response {
	req.Response = rs.getOrCreateResponse(req)
	stored, err := rs.RequestRepository.Add(req)

	if err != nil {
		req.Response.Error = err.Error()
		return req.Response
	}

	if rs.Hub != nil {
		rs.Hub.Publish(RequestEvent{
			Path:      stored.Path,
			ID:        stored.ID,
			Method:    stored.Method,
			Status:    stored.Response.ResponseStatusCode,
			Timestamp: time.Unix(0, stored.Timestamp).Format("2 Jan 2006 15:04:05"),
			UnixMs:    stored.Timestamp / int64(time.Millisecond),
		})
	}

	return req.Response
}

func (rs RequestService) getOrCreateResponse(req request.Request) request.Response {
	allConfigs := rs.ConfigRepository.FindAll()

	configItem := config.FindMatchingConfig(allConfigs, req.Path)
	if configItem == nil {
		return request.Response{
			Target:             "no target - mocked response",
			ResponseStatusCode: defaultStatusCode,
			ResponseBody:       "ok",
		}
	}

	if configItem.HasTargetPath() {
		targetPath := configItem.GenerateFinalTargetPath(req.Path)
		response := rs.Transport.SendRequestToTarget(req, targetPath)
		response.MatchedConfig = configItem.Path
		return response
	} else {
		return request.Response{
			Target:             "no target - mocked response",
			MatchedConfig:      configItem.Path,
			ResponseStatusCode: configItem.Response.StatusCode,
			ResponseBody:       configItem.Response.Body,
		}
	}
}
