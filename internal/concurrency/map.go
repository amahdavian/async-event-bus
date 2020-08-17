package concurrency

import (
	"sync"
)

type MapOfSlice struct {
	items map[string][]interface{}
	mutex sync.RWMutex
}

func NewMapOfSlice() *MapOfSlice {
	return &MapOfSlice{items: map[string][]interface{}{}}
}

func (concurrentMap *MapOfSlice) Get(key string) ([]interface{}, bool) {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()
	if items, found := concurrentMap.items[key]; found {
		copyOfItems := append([]interface{}{}, items...)
		return copyOfItems, true
	}
	return nil, false
}

func (concurrentMap *MapOfSlice) AppendAt(key string, item interface{}) {
	concurrentMap.mutex.Lock()
	defer concurrentMap.mutex.Unlock()
	concurrentMap.items[key] = append(concurrentMap.items[key], item)
}

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