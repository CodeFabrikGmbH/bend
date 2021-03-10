package request

type Transport interface {
	SendRequestToTarget(request Request, targetUrl string) Response
}
