package validator

import (
    "bytes"
    "encoding/json"
    "net/http"

    "github.com/janospapp/json-validator/app"
    "github.com/janospapp/json-validator/schema"
    "github.com/santhosh-tekuri/jsonschema/v5"
)

func Check(id string, doc []byte, store schema.Store) ([]byte, int) {
    resp := app.Response{
        Action: app.ACTION_VALIDATE,
        Id: id,
        Status: app.SUCCESS,
    }
    code := http.StatusOK

    var docCheck, schemaFound bool
    var sch []byte

    if schema.CheckId(id, &resp, &code) {
        docCheck = app.CheckJSON(doc, &resp, &code)
    }

    if docCheck {
        sch, schemaFound = store.GetSchema(id)
        if schemaFound {
            validate(id, doc, sch, &resp)
        } else {
            resp.Status = app.ERROR
            resp.Message = "Schema not found"
            code = http.StatusNotFound
        }
    }

    return binary(resp), code
}

func validate(id string, doc []byte, schema []byte, resp *app.Response) {
    sch, err := jsonschema.CompileString(id, bytes.NewBuffer(schema).String())
    var v interface{}
    json.Unmarshal(doc, &v)
    cleanNulls([]*interface{}{&v})
    if err = sch.Validate(v); err != nil {
        resp.Status = app.ERROR
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

func binary(resp app.Response) []byte {
    json, _ := json.MarshalIndent(resp, "", "  ")
    return json
}
