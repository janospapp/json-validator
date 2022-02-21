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

    switch {
    case CheckId(id, &resp, &code) == false:
    case app.CheckJSON(schema, &resp, &code) == false:
    case store.StoreSchema(id, schema) == false:
        resp.Status = app.ERROR
        resp.Message = "Couldn't save your schema. Please contact support."
        code = http.StatusInternalServerError
    }

    return resp.Bytes(), code
}

func Get(id string, store Store) ([]byte, int) {
    resp := app.Response{
        Action: app.ACTION_GET,
        Id: id,
        Status: app.ERROR,
    }
    var code int

    switch {
    case CheckId(id, &resp, &code) == false:
    default:
        schema, found := store.GetSchema(id)
        if !found {
            resp.Message = "Schema not found"
            code = http.StatusNotFound
        } else {
            // Indenting the schema is not necessary, it only
            // makes the output more readable
            return indent(schema), http.StatusOK
        }
    }

    return resp.Bytes(), code
}

func indent(schema []byte) []byte {
    buf := bytes.NewBuffer(make([]byte, len(schema)))
    json.Indent(buf, schema, "", "  ")

    return buf.Bytes()
}
