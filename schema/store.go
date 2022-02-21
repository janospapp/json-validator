package schema

type Store interface {
    StoreSchema(string, []byte) bool
    GetSchema(string) ([]byte, bool)
}

type MemoryStore struct {
    Schemas map[string][]byte
}

func NewMemoryStore() *MemoryStore {
    store := MemoryStore{
        Schemas: make(map[string][]byte),
    }

    return &store
}

func (store *MemoryStore) StoreSchema(id string, schema []byte) bool {
    store.Schemas[id] = schema
    return true
}

func (store *MemoryStore) GetSchema(id string) ([]byte, bool) {
    schema, ok := store.Schemas[id]
    return schema, ok
}
