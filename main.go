package main

import (
    "io"
    "log"
    "net/http"

    "github.com/janospapp/json-validator/schema"
)

func main() {
    http.HandleFunc("/schema/", schemaHandler)
    log.Fatal(http.ListenAndServe(":8000", nil))
}

func schemaHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/schema/"):]
    var response []byte
    var code int

    switch r.Method {
    case http.MethodPost:
        // Note: ReadAll is used for simplicity. A production
        // application shouldn't use it, due to security reasons.
        // Content-Length must be checked before reading to filter
        // out too long requests, or read the body chunk by chunk.
        // This would be covered by a proper web framework.
        body, err := io.ReadAll(r.Body)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        response, code = schema.Upload(id, body)
    case http.MethodGet:
        response, code = schema.Get(id)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    if code != 0 {
        w.WriteHeader(code)
    }

    w.Write(response)
}
