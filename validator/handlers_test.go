package validator

import (
    "encoding/json"
    "net/http"
    "testing"

    "github.com/janospapp/json-validator/schema"
)

var sch = []byte(`{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "type": "object",
    "properties": {
        "name": {
            "type": "string"
        },
        "age": {
            "type": "integer",
            "minimum": 0,
            "maximum": 99
        },
        "job": {
            "type": "string"
        }
    },
    "required": ["name", "age"]
}`)

func TestCheckEmptySchemaId(t *testing.T) {
    id := ""
    doc := []byte{}
    store := schema.NewMemoryStore()

    resp, code := Check(id, doc, store)
    if r := getResp(resp); r.Status != ERROR {
        t.Fatal("ValidateDoc with empty schema id must be an error")
    }

    if code != http.StatusBadRequest {
        t.Fatalf("ValidateDoc with empty schema id status code was %d, expected is 400", code)
    }
}

func TestCheckInvalidDocument(t *testing.T) {
    id := "test"
    doc := []byte(`{"document": [1, 2, 3`)
    store := schema.NewMemoryStore()

    resp, code := Check(id, doc, store)
    if r := getResp(resp); r.Status != ERROR {
        t.Fatal("Validating invalid JSON document must be an error")
    }

    if code != http.StatusBadRequest {
        t.Fatalf("Validating invalid JSON status code was %d, expected is 400", code)
    }
}

func TestCheckNonExistentSchema(t *testing.T) {
    id := "test"
    doc := []byte(`{"user": "admin"}`)
    store := schema.NewMemoryStore()

    resp, code := Check(id, doc, store)
    if r := getResp(resp); r.Status != ERROR {
        t.Fatal("Validating against non existent schema must be an error")
    }

    if code != http.StatusNotFound {
        t.Fatalf("Validating against non existent schema status code was %d, expected is 404", code)
    }
}

func TestCheckSuccess(t *testing.T) {
    id := "test"
    doc := []byte(`{
        "name": "Steve",
        "age": 43,
        "job": "doctor"
    }`)
    store := schema.NewMemoryStore()
    store.StoreSchema(id, sch)

    resp, code := Check(id, doc, store)
    if r := getResp(resp); r.Status != SUCCESS {
        t.Fatalf("ValidateDoc is expected to succeed. Error: %s", r.Message)
    }

    if code != http.StatusOK {
        t.Fatalf("ValidateDoc status code was %d, expected is 200", code)
    }
}

func TestCheckRemovesNullValues(t *testing.T) {
    id := "test"
    doc := []byte(`{
        "name": "Steve",
        "age": 43,
        "hobbies": null
    }`)
    store := schema.NewMemoryStore()
    store.StoreSchema(id, sch)

    resp, code := Check(id, doc, store)
    if r := getResp(resp); r.Status != SUCCESS {
        t.Fatalf("ValidateDoc with null value is expected to succeed. Error: %s", r.Message)
    }

    if code != http.StatusOK {
        t.Fatalf("ValidateDoc with null value status code was %d, expected is 200", code)
    }
}

func getResp(data []byte) Response {
    var resp Response
    json.Unmarshal(data, &resp)
    return resp
}
