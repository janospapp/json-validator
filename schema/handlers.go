package schema

import (
    "encoding/json"
    "net/http"
)

const (
    SUCCESS = "success"
    ERROR = "error"
    ACTION_UPLOAD = "uploadSchema"
    ACTION_GET = "getSchema"
)

type Response struct {
    Action  string `json:"action"`
    Id      string `json:"id"`
    Status  string `json:"status"`
    Message string `json:"message"`
}

func Upload(id string, schema []byte) ([]byte, int) {
    r := Response{
        Action: ACTION_UPLOAD,
        Id: id,
        Status: ERROR,
        Message: "Not implemented",
    }

    json, _ := json.MarshalIndent(r, "", "  ")
    return json, http.StatusMethodNotAllowed
}

func Get(id string) ([]byte, int) {
    r := Response{
        Action: ACTION_GET,
        Id: id,
        Status: ERROR,
        Message: "Not implemented",
    }

    json, _ := json.MarshalIndent(r, "", "  ")
    return json, http.StatusMethodNotAllowed
}
