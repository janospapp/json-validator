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
    Message string `json:"message,omitempty"`
}

func Upload(id string, schema []byte) ([]byte, int) {
    resp := Response{
        Action: ACTION_UPLOAD,
        Id: id,
        Status: SUCCESS,
    }
    code := http.StatusCreated

    var schemaCheck, stored bool
    idCheck := checkId(id, &resp, &code)

    if idCheck {
        schemaCheck = checkSchema(schema, &resp, &code)
    }

    if schemaCheck {
        stored = storeSchema(id, schema)
    }

    if schemaCheck && !stored {
        resp.Status = ERROR
        resp.Message = "Couldn't save your schema. Please contact support."
        code = http.StatusInternalServerError
    }

    json, _ := json.MarshalIndent(resp, "", "  ")
    return json, code
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

func checkId(id string, resp *Response, code *int) bool {
    if id == "" {
        resp.Status = ERROR
        resp.Message = "id cannot be empty"
        *code = http.StatusBadRequest
        return false
    }

    return true
}

func checkSchema(schema []byte, resp *Response, code *int) bool {
    if !json.Valid(schema) {
        resp.Status = ERROR
        resp.Message = "Invalid JSON"
        *code = http.StatusBadRequest
        return false
    }

    return true
}
