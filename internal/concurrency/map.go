package concurrency

import (
	"sync"
)

// MapOfSlice maps keys of type string to values of type []interface{} and allows concurrent reads and non-concurrent writes.
type MapOfSlice struct {
	items map[string][]interface{}
	mutex sync.RWMutex
}

// NewMapOfSlice creates a new instance of MapOfSlice
func NewMapOfSlice() *MapOfSlice {
	return &MapOfSlice{items: map[string][]interface{}{}}
}

// Get returns the value to which the specified key is mapped, or false if this map contains no mapping for the key
func (concurrentMap *MapOfSlice) Get(key string) ([]interface{}, bool) {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()
	if items, found := concurrentMap.items[key]; found {
		copyOfItems := append([]interface{}{}, items...)
		return copyOfItems, true
	}
	return nil, false
}

// AppendAt appends the specified value to the end of the slice associated with the specified key in this map
func (concurrentMap *MapOfSlice) AppendAt(key string, item interface{}) {
	concurrentMap.mutex.Lock()
	defer concurrentMap.mutex.Unlock()
	concurrentMap.items[key] = append(concurrentMap.items[key], item)
}

// RemoveAt removes the specified value from the slice associated with the specified key in this map
func (concurrentMap *MapOfSlice) RemoveAt(key string, item interface{}) {
	concurrentMap.mutex.Lock()
	defer concurrentMap.mutex.Unlock()
	if items, found := concurrentMap.items[key]; found {
		concurrentMap.items[key] = remove(items, item)
	}
}

func remove(items []interface{}, item interface{}) []interface{} {
	if index, found := find(items, item); found {
		return append(items[:index], items[index+1:]...)
	}
	return items
}

func find(items []interface{}, itemToFind interface{}) (int, bool) {
	for index, item := range items {
		if itemToFind == item {
			return index, true
		}
	}
	return len(items), false
}
