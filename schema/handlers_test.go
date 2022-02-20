package schema

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"
)

type BrokenStore struct {}

func (BrokenStore) storeSchema(string, []byte) bool {
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
    if r := getResp(resp); r.Status != ERROR {
        t.Fatal("Uploading empty id schema must be treated as errors")
    }

    if code != http.StatusBadRequest {
        t.Fatalf("Uploading empty id schema status code was %d, expected is 400", code)
    }
}

func TestUploadInvalidJSONSchema(t *testing.T) {
    id := "test"
    data := []byte("{bad: json")

    resp, code := Upload(id, data, BrokenStore{})
    if r := getResp(resp); r.Status != ERROR {
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
    if r := getResp(resp); r.Status != ERROR {
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
    if r := getResp(resp); r.Status != SUCCESS {
        t.Fatalf("Upload failed with message: %s. Success was expected", r.Message)
    }

    if code != http.StatusCreated {
        t.Fatalf("Upload schema status code was %d, expected is 201", code)
    }
}

func TestGetEmptyIdSchema(t *testing.T) {
    id := ""

    resp, code := Get(id, BrokenStore{})
    if r := getResp(resp); r.Status != ERROR {
        t.Fatal("Getting empty id schema must be treated as error")
    }

    if code != http.StatusBadRequest {
        t.Fatalf("Get empty id schema status code was %d, expected is 400", code)
    }
}

func TestGetNonExistenSchema(t *testing.T) {
    id := "sample"

    resp, code := Get(id, BrokenStore{})
    if r := getResp(resp); r.Status != ERROR {
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
    store.storeSchema(id, data)

    resp, code := Get(id, store)
    if r := bytes.Compare(resp, data); r != 0 {
        t.Fatalf("Getting schema returned %s, empty document is expected", resp)
    }

    if code != http.StatusOK {
        t.Fatalf("Get schema status code was %d, expected is 200", code)
    }
}

func getResp(data []byte) Response {
    var resp Response
    json.Unmarshal(data, &resp)
    return resp
}
