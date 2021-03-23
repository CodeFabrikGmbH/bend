package application

import (
	"code-fabrik.com/bend/domain/config"
	"code-fabrik.com/bend/domain/request"
	"net/http"
)

type RequestService struct {
	RequestRepository request.Repository
	ConfigRepository  config.Repository
	Transport         request.Transport
}

func (rs RequestService) Delete(path string, requestId *string) error {
	if requestId == nil {
		return rs.RequestRepository.DeletePath(path)
	} else {
		return rs.RequestRepository.DeleteRequestForPath(path, *requestId)
	}
}

func (rs RequestService) SendRequestToTarget(path, requestId, targetUrl string) request.Response {
	req := rs.RequestRepository.GetRequest(path, requestId)
	return rs.Transport.SendRequestToTarget(req, targetUrl)
}

func (rs RequestService) HandleTrackableRequest(req request.Request) request.Response {
	req.Response = rs.getOrCreateResponse(req)
	err := rs.RequestRepository.Add(req)

	if err != nil {
		req.Response.Error = err.Error()
	}

	return req.Response
}

func (rs RequestService) getOrCreateResponse(req request.Request) request.Response {
	config := rs.ConfigRepository.Find(req.Path)

	if config == nil {
		return request.Response{
			Target:             "no target - mocked response",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       "ok",
		}
	}
	if len(config.Target) != 0 {
		return rs.Transport.SendRequestToTarget(req, config.Target)
	}

	return request.Response{
		Target:             "no target - mocked response",
		ResponseStatusCode: config.Response.StatusCode,
		ResponseBody:       config.Response.Body,
	}
}
