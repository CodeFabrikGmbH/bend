package request

type Response struct {
	Error              string              `json:"error,omitempty"`
	Target             string              `json:"target,omitempty"`
	MatchedConfig      string              `json:"matchedConfig,omitempty"` // path of the endpoint config that handled the request
	ResponseBody       string              `json:"responseBody,omitempty"`
	ResponseHeader     map[string][]string `json:"responseHeader,omitempty"`
	ResponseStatusCode int                 `json:"responseStatusCode,omitempty"`
}
