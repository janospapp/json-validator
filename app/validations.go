package app

import (
    "encoding/json"
    "net/http"
)

func CheckJSON(doc []byte, resp *Response, code *int) bool {
    if !json.Valid(doc) {
        resp.Status = ERROR
        resp.Message = "Invalid JSON"
        *code = http.StatusBadRequest
        return false
    }

    return true
}

