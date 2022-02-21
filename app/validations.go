package app

import (
    "encoding/json"
    "log"
    "net/http"
)

func CheckJSON(doc []byte, resp *Response, code *int) bool {
    if !json.Valid(doc) {
        log.Printf("JSON is not valid: \n%s\n", doc)
        resp.Status = ERROR
        resp.Message = "Invalid JSON"
        *code = http.StatusBadRequest
        return false
    }

    return true
}

