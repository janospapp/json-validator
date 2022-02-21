package schema

import (
    "errors"
    "os"
)

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

type FileStore struct {
    BaseDir string
}

func NewFileStore() *FileStore {
    store := FileStore{
        BaseDir: "saved_schemas/",
    }

    return &store
}

func (store *FileStore) StoreSchema(id string, schema []byte) bool {
    // Ensure that the directory exists
    if _, err := os.Stat(store.BaseDir); errors.Is(err, os.ErrNotExist) {
        os.Mkdir(store.BaseDir, 0744)
    }

    err := os.WriteFile(store.getPath(id), schema, 0644)
    return err == nil
}

func (store *FileStore) GetSchema(id string) ([]byte, bool) {
    schema, err := os.ReadFile(store.getPath(id))
    return schema, err == nil
}

func (store *FileStore) getPath(name string) string {
    return store.BaseDir + name + ".json"
}
