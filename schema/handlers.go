package schema

import (
    "bytes"
    "encoding/json"
    "net/http"

    "github.com/janospapp/json-validator/app"
)

func Upload(id string, schema []byte, store Store) ([]byte, int) {
    resp := app.Response{
        Action: app.ACTION_UPLOAD,
        Id: id,
        Status: app.SUCCESS,
    }
    code := http.StatusCreated

    var schemaCheck, stored bool
    idCheck := CheckId(id, &resp, &code)

    if idCheck {
        schemaCheck = app.CheckJSON(schema, &resp, &code)
    }

    if schemaCheck {
        stored = store.StoreSchema(id, schema)
    }

    if schemaCheck && !stored {
        resp.Status = app.ERROR
        resp.Message = "Couldn't save your schema. Please contact support."
        code = http.StatusInternalServerError
    }

    return binary(resp), code
}

func Get(id string, store Store) ([]byte, int) {
    resp := app.Response{
        Action: app.ACTION_GET,
        Id: id,
        Status: app.ERROR,
    }
    var code int

    if CheckId(id, &resp, &code) {
        schema, found := store.GetSchema(id)
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

func binary(resp app.Response) []byte {
    json, _ := json.MarshalIndent(resp, "", "  ")
    return json
}

func indent(schema []byte) []byte {
    buf := bytes.NewBuffer(make([]byte, len(schema)))
    json.Indent(buf, schema, "", "  ")

    return buf.Bytes()
}
