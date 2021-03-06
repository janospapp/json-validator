package app

import "encoding/json"

const (
    SUCCESS = "success"
    ERROR = "error"
    ACTION_UPLOAD = "uploadSchema"
    ACTION_GET = "getSchema"
    ACTION_VALIDATE = "validateDocument"
)

type Response struct {
    Action  string `json:"action"`
    Id      string `json:"id"`
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`
}

func (resp Response) Bytes() []byte {
    json, _ := json.MarshalIndent(resp, "", "  ")
    return json
}
