package schema

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"

    "github.com/janospapp/json-validator/app"
)

type BrokenStore struct {}

func (BrokenStore) StoreSchema(string, []byte) bool {
    // Always report unsuccessful save
    return false
}

func (BrokenStore) GetSchema(string) ([]byte, bool) {
    // Always return not found
    return nil, false
}

func TestUploadEmptyIdSchema(t *testing.T) {
    id := ""
    data := []byte("")

    resp, code := Upload(id, data, BrokenStore{})
    if r := getResp(resp); r.Status != app.ERROR {
        t.Fatal("Uploading empty id schema must be treated as errors")
    }

    if code != http.StatusBadRequest {
        t.Fatalf("Uploading empty id schema status code was %d, expected is 400", code)
    }
}

func TestUploadExistingId(t *testing.T) {
    id := "saved"
    data := []byte("")
    store := NewMemoryStore()
    store.StoreSchema(id, data)

    resp, code := Upload(id, data, store)
    if r := getResp(resp); r.Status != app.ERROR {
        t.Fatal("Uploading not unique schema must be treated as errors")
    }

    if code != http.StatusConflict {
        t.Fatalf("Uploading empty id schema status code was %d, expected is 409", code)
    }
}

func TestUploadInvalidJSONSchema(t *testing.T) {
    id := "test"
    data := []byte("{bad: json")

    resp, code := Upload(id, data, BrokenStore{})
    if r := getResp(resp); r.Status != app.ERROR {
        t.Fatal("Uploading invalid JSON as schema must be treated as error")
    }

    if code != http.StatusBadRequest {
        t.Fatalf("Uploading invalid JSON schema status code was %d, expected is 400", code)
    }
}

func TestUploadFailedToStoreSchema(t *testing.T) {
    id := "test"
    data := []byte(`{"type": "object", "properties": []}`)

    resp, code := Upload(id, data, BrokenStore{})
    if r := getResp(resp); r.Status != app.ERROR {
        t.Fatal("Upload failing to store the schema must be treated as error")
    }

    if code != http.StatusInternalServerError {
        t.Fatalf("Uploading invalid JSON schema status code was %d, expected is 500", code)
    }
}

func TestUploadSchemaSuccess(t *testing.T) {
    id := "test"
    data := []byte(`{"type": "object", "properties": []}`)

    store := NewMemoryStore()
    resp, code := Upload(id, data, store)
    if r := getResp(resp); r.Status != app.SUCCESS {
        t.Fatalf("Upload failed with message: %s. Success was expected", r.Message)
    }

    if code != http.StatusCreated {
        t.Fatalf("Upload schema status code was %d, expected is 201", code)
    }
}

func TestGetEmptyIdSchema(t *testing.T) {
    id := ""

    resp, code := Get(id, BrokenStore{})
    if r := getResp(resp); r.Status != app.ERROR {
        t.Fatal("Getting empty id schema must be treated as error")
    }

    if code != http.StatusBadRequest {
        t.Fatalf("Get empty id schema status code was %d, expected is 400", code)
    }
}

func TestGetNonExistenSchema(t *testing.T) {
    id := "sample"

    resp, code := Get(id, BrokenStore{})
    if r := getResp(resp); r.Status != app.ERROR {
        t.Fatal("Getting non existent schema must be treated as error")
    }

    if code != http.StatusNotFound {
        t.Fatalf("Get non existing schema status code was %d, expected is 404", code)
    }
}

func TestGetSchemaSuccess(t *testing.T) {
    id := "test"
    data := []byte{}

    store := NewMemoryStore()
    store.StoreSchema(id, data)

    resp, code := Get(id, store)
    if r := bytes.Compare(resp, data); r != 0 {
        t.Fatalf("Getting schema returned %s, empty document is expected", resp)
    }

    if code != http.StatusOK {
        t.Fatalf("Get schema status code was %d, expected is 200", code)
    }
}

func getResp(data []byte) app.Response {
    var resp app.Response
    json.Unmarshal(data, &resp)
    return resp
}
