package main

import (
    "io"
    "log"
    "net/http"

    "github.com/janospapp/json-validator/schema"
    "github.com/janospapp/json-validator/validator"
)

var store schema.Store

func main() {
    log.SetPrefix("json-validator: ")
    log.SetFlags(0)
    // A production application would store the schemas
    // in a robust database, instead of directly on the
    // filesystem. This is only for persistency in the
    // exercise.
    store = schema.NewFileStore()

    http.HandleFunc("/schema/", schemaHandler)
    http.HandleFunc("/validate/", validateHandler)

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
            log.Printf("Couldn't read request body due to '%s'\n", err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        response, code = schema.Upload(id, body, store)
    case http.MethodGet:
        response, code = schema.Get(id, store)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    w.WriteHeader(code)
    w.Write(response)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/validate/"):]
    var response []byte
    var code int

    switch r.Method {
    case http.MethodPost:
        // See the comment about ReadAll above
        body, err := io.ReadAll(r.Body)
        if err != nil {
            log.Printf("Couldn't read request body due to '%s'\n", err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        response, code = validator.Check(id, body, store)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    w.WriteHeader(code)
    w.Write(response)
}
