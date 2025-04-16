package keyValueDb

import (
	"fmt"
	"imageresizerservice/library/uow"
)

// ImplNamespaced implements the KeyValueDb interface by wrapping another KeyValueDb
// and prefixing all keys with a namespace
type ImplNamespaced struct {
	db        KeyValueDb
	namespace string
}

// NewImplNamespaced creates a new instance of ImplNamespaced
func NewImplNamespaced(db KeyValueDb, namespace string) *ImplNamespaced {
	return &ImplNamespaced{
		db:        db,
		namespace: namespace,
	}
}

// namespaceKey prefixes the key with the namespace
func (db *ImplNamespaced) namespaceKey(key string) string {
	return fmt.Sprintf("%s:%s", db.namespace, key)
}

// Get retrieves a value by key. Returns nil if key not found.
func (db *ImplNamespaced) Get(key string) (*string, error) {
	namespacedKey := db.namespaceKey(key)
	return db.db.Get(namespacedKey)
}

// Put stores a key-value pair
func (db *ImplNamespaced) Put(uow *uow.Uow, key string, value string) error {
	namespacedKey := db.namespaceKey(key)
	return db.db.Put(uow, namespacedKey, value)
}

// Zap removes a key-value pair
func (db *ImplNamespaced) Zap(uow *uow.Uow, key string) error {
	namespacedKey := db.namespaceKey(key)
	return db.db.Zap(uow, namespacedKey)
}

// Ensure ImplNamespaced implements KeyValueDb interface
var _ KeyValueDb = (*ImplNamespaced)(nil)
