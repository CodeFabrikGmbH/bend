package application

import (
	"code-fabrik.com/bend/domain/environment"
	"code-fabrik.com/bend/domain/request"
	"net/http"
)

type RequestService struct {
	Env environment.Environment
}

func (rs RequestService) Delete(path string, requestId *string) error {
	if requestId == nil {
		return rs.Env.RequestRepository.DeletePath(path)
	} else {
		return rs.Env.RequestRepository.DeleteRequestForPath(path, *requestId)
	}
}

func (rs RequestService) SendRequestToTarget(path, requestId, targetUrl string) request.Response {
	req := rs.Env.RequestRepository.GetRequest(path, requestId)
	return rs.Env.Transport.SendRequestToTarget(req, targetUrl)
}

func (rs RequestService) HandleTrackableRequest(req request.Request) request.Response {
	req.Response = rs.getOrCreateResponse(req)
	err := rs.Env.RequestRepository.Add(req)

	if err != nil {
		req.Response.Error = err.Error()
	}

	return req.Response
}

func (rs RequestService) getOrCreateResponse(req request.Request) request.Response {
	config := rs.Env.ConfigRepository.Find(req.Path)

	if config == nil {
		return request.Response{
			Target:             "no target - mocked response",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       "ok",
		}
	}
	if len(config.Target) != 0 {
		return rs.Env.Transport.SendRequestToTarget(req, config.Target)
	}

	return request.Response{
		Target:             "no target - mocked response",
		ResponseStatusCode: config.Response.StatusCode,
		ResponseBody:       config.Response.Body,
	}
}
