package request

type Response struct {
	Error              string `json:"error,omitempty"`
	ResponseBody       string `json:"responseBody,omitempty"`
	ResponseStatusCode int    `json:"responseStatusCode,omitempty"`
}
