package application

import (
	"code-fabrik.com/bend/domain/config"
	"code-fabrik.com/bend/domain/request"
)

type RequestService struct {
	RequestRepository request.Repository
	ConfigRepository  config.Repository
	Transport         request.Transport
}

const defaultStatusCode = 200

func (rs RequestService) DeleteAllRequestsForPath(path string) error {
	return rs.RequestRepository.DeletePath(path)
}

func (rs RequestService) DeleteRequest(path string, requestId string) error {
	return rs.RequestRepository.DeleteRequestForPath(path, requestId)
}

func (rs RequestService) SendRequestToTarget(path, requestId, targetUrl string) request.Response {
	req := rs.RequestRepository.GetRequest(path, requestId)
	return rs.Transport.SendRequestToTarget(req, targetUrl)
}

func (rs RequestService) TrackRequest(req request.Request) request.Response {
	req.Response = rs.getOrCreateResponse(req)
	err := rs.RequestRepository.Add(req)

	if err != nil {
		req.Response.Error = err.Error()
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
		return rs.Transport.SendRequestToTarget(req, targetPath)
	} else {
		return request.Response{
			Target:             "no target - mocked response",
			ResponseStatusCode: configItem.Response.StatusCode,
			ResponseBody:       configItem.Response.Body,
		}
	}
}
