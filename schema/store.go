package schema

var schemas = make(map[string][]byte)

func storeSchema(id string, schema []byte) bool {
    schemas[id] = schema
    return true
}

func GetSchema(id string) ([]byte, bool) {
    schema, ok := schemas[id]
    return schema, ok
}
