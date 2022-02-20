package validator

import (
    "bytes"
    "encoding/json"
    "net/http"

    "github.com/janospapp/json-validator/schema"
    "github.com/santhosh-tekuri/jsonschema/v5"
)

const (
    SUCCESS = "success"
    ERROR = "error"
    ACTION = "validateDocument"
)

type Response struct {
    Action  string `json:"action"`
    Id      string `json:"id"`
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`
}

func Check(id string, doc []byte) ([]byte, int) {
    resp := Response{
        Action: ACTION,
        Id: id,
        Status: SUCCESS,
    }
    code := http.StatusOK

    var docCheck, schemaFound bool
    var sch []byte

    if checkId(id, &resp, &code) {
        docCheck = checkDoc(doc, &resp, &code)
    }

    if docCheck {
        sch, schemaFound = schema.GetSchema(id)
    }

    if schemaFound {
        validate(id, doc, sch, &resp)
    } else {
        resp.Status = ERROR
        resp.Message = "Schema not found"
        code = http.StatusNotFound
    }

    return binary(resp), code
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

func checkDoc(doc []byte, resp *Response, code *int) bool {
    if !json.Valid(doc) {
        resp.Status = ERROR
        resp.Message = "Input document is invalid JSON"
        *code = http.StatusBadRequest
        return false
    }

    return true
}

func validate(id string, doc []byte, schema []byte, resp *Response) {
    sch, err := jsonschema.CompileString(id, bytes.NewBuffer(schema).String())
    var v interface{}
    json.Unmarshal(doc, &v)
    if err = sch.Validate(v); err != nil {
        resp.Status = ERROR
        resp.Message = fetchTheFirstError(err.(*jsonschema.ValidationError))
    }
}

func fetchTheFirstError(err *jsonschema.ValidationError) string {
    if len(err.Causes) > 0 {
        return fetchTheFirstError(err.Causes[0])
    } else {
        return "/root" + err.InstanceLocation + ": " + err.Message
    }
}

func binary(resp Response) []byte {
    json, _ := json.MarshalIndent(resp, "", "  ")
    return json
}
