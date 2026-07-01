package request

import "encoding/json"

type Request struct {
	ID        string              `json:"ID"`
	Timestamp int64               `json:"timestamp"`
	Path      string              `json:"path"`
	Method    string              `json:"method"`
	Body      string              `json:"body"`
	Header    map[string][]string `json:"header"`
	Host      string              `json:"host"`
	Uri       string              `json:"uri"`
	Response  Response            `json:"response"`
}

func (r Request) ToJson() string {
	marshal, _ := json.Marshal(r)
	return string(marshal)
}

// Summary is a lightweight projection of a Request used for list views. It
// deliberately omits the body, headers and response so those large fields are
// not read into memory when only an overview is needed.
type Summary struct {
	ID        string `json:"ID"`
	Timestamp int64  `json:"timestamp"`
}
