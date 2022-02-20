package schema

import (
    "bytes"
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

    return binary(resp), code
}

func Get(id string) ([]byte, int) {
    resp := Response{
        Action: ACTION_GET,
        Id: id,
        Status: ERROR,
    }
    var code int

    if checkId(id, &resp, &code) {
        schema, found := GetSchema(id)
        if !found {
            resp.Message = "Schema not found"
            return binary(resp), http.StatusNotFound
        } else {
            // Indenting the schema is not necessary, it only
            // makes the output more readable
            return indent(schema), http.StatusOK
        }
    } else {
        return binary(resp), code
    }
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

func binary(resp Response) []byte {
    json, _ := json.MarshalIndent(resp, "", "  ")
    return json
}

func indent(schema []byte) []byte {
    buf := bytes.NewBuffer(make([]byte, len(schema)))
    json.Indent(buf, schema, "", "  ")

    return buf.Bytes()
}
