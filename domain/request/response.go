package request

type Response struct {
	Error              string `json:"error,omitempty"`
	Target             string `json:"target,omitempty"`
	ResponseBody       string `json:"responseBody,omitempty"`
	ResponseStatusCode int    `json:"responseStatusCode,omitempty"`
}
