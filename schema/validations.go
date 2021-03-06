package schema

import (
    "net/http"

    "github.com/janospapp/json-validator/app"
)

func CheckId(id string, resp *app.Response, code *int) bool {
    if id == "" {
        resp.Status = app.ERROR
        resp.Message = "id cannot be empty"
        *code = http.StatusBadRequest
        return false
    }

    return true
}

func UniqueId(id string, store Store, resp *app.Response, code *int) bool {
    exists := store.Exists(id)
    if exists {
        resp.Status = app.ERROR
        resp.Message = "Schema already exists"
        *code = http.StatusConflict
    }

    return !exists
}
