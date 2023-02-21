package server

import (
	"sync"

	"github.com/google/uuid"
)

// dataType is an abstraction of the internal type of the storage
type dataType string

func (t dataType) string() string {
	return string(t)
}

// DataStore is and intermediate storage layer between sync providers and server stream handler
type DataStore struct {
	data          dataType
	subscriptions map[string]chan dataType

	mu sync.RWMutex
}

func NewDataStore() *DataStore {
	return &DataStore{
		data:          "",
		subscriptions: make(map[string]chan dataType),
	}
}

func (store *DataStore) subscribe(id string, c chan dataType) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.subscriptions[id] = c
}

func (store *DataStore) unsubscribe(id string) {
	store.mu.Lock()
	defer store.mu.Unlock()

	delete(store.subscriptions, id)
}

func (store *DataStore) cache(data dataType) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.data = data

	for _, sub := range store.subscriptions {
		sub <- data
	}
}

func (store *DataStore) currentState() dataType {
	store.mu.RLock()
	defer store.mu.RUnlock()

	return store.data
}

// StorageID is an abstraction to generate unique storage subscription identifiers
func StorageID() string {
	return uuid.New().String()
}
