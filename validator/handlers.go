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
    cleanNulls([]*interface{}{&v})
    if err = sch.Validate(v); err != nil {
        resp.Status = ERROR
        resp.Message = fetchTheFirstError(err.(*jsonschema.ValidationError))
    }
}

// This function checks the JSON document level by level starting
// at the top. If a value is a nested JSON object, then its nil
// values are deleted and all its other values are added for checks
// at the next level. It ends when there are no more levels to check.
func cleanNulls(values []*interface{}) {
    var rest []*interface{}
    for _, obj := range values {
        switch (*obj).(type) {
        case map[string]interface{}:
            // The current value is a JSON object, check its values
            m := (*obj).(map[string]interface{})
            for k, v := range m {
                if v == nil {
                    delete(m, k)
                } else {
                    // Add the value for further checking
                    rest = append(rest, &v)
                }
            }
        }
    }

    if len(rest) > 0 {
        cleanNulls(rest)
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
