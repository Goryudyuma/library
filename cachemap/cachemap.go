package cachemap

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"maps"
	"sync"
)

type CacheMap[K comparable, V any] struct {
	m    map[K]V
	lock sync.RWMutex
}

func NewCacheMap[K comparable, V any]() *CacheMap[K, V] {
	cm := CacheMap[K, V]{m: make(map[K]V)}
	return &cm
}

func (cm *CacheMap[K, V]) Lock() {
	cm.lock.Lock()
}

func (cm *CacheMap[K, V]) Unlock() {
	cm.lock.Unlock()
}

func (cm *CacheMap[K, V]) RLock() {
	cm.lock.RLock()
}

func (cm *CacheMap[K, V]) RUnlock() {
	cm.lock.RUnlock()
}

func (cm *CacheMap[K, V]) Delete(key K) {
	cm.Lock()
	defer cm.Unlock()
	cm.DeleteWithoutLock(key)
}

func (cm *CacheMap[K, V]) DeleteWithoutLock(key K) {
	delete(cm.m, key)
}

func (cm *CacheMap[K, V]) Get(key K) (V, bool) {
	cm.RLock()
	defer cm.RUnlock()
	return cm.GetWithoutLock(key)
}

func (cm *CacheMap[K, V]) GetWithoutLock(key K) (value V, ok bool) {
	value, ok = cm.m[key]
	return
}

func (cm *CacheMap[K, V]) Set(key K, value V) {
	cm.Lock()
	defer cm.Unlock()
	cm.SetWithoutLock(key, value)
}

func (cm *CacheMap[K, V]) SetWithoutLock(key K, value V) {
	cm.m[key] = value
}

func (cm *CacheMap[K, V]) Len() int {
	cm.RLock()
	defer cm.RUnlock()
	return cm.LenWithoutLock()
}

func (cm *CacheMap[K, V]) LenWithoutLock() int {
	return len(cm.m)
}

func (cm *CacheMap[K, V]) Clear() {
	cm.Lock()
	defer cm.Unlock()
	cm.ClearWithoutLock()
}

func (cm *CacheMap[K, V]) ClearWithoutLock() {
	cm.m = make(map[K]V)
}

func (cm *CacheMap[K, V]) Export() map[K]V {
	cm.RLock()
	defer cm.RUnlock()
	return cm.ExportWithoutLock()
}

func (cm *CacheMap[K, V]) ExportWithoutLock() map[K]V {
	return maps.Collect(maps.All(cm.m))
}

func (cm *CacheMap[K, V]) Import(m map[K]V) {
	cm.Lock()
	defer cm.Unlock()
	cm.ImportWithoutLock(m)
}

func (cm *CacheMap[K, V]) ImportWithoutLock(m map[K]V) {
	cm.m = m
}

func (cm *CacheMap[K, V]) GobEncode() ([]byte, error) {
	cm.RLock()
	defer cm.RUnlock()
	// cm.mをGobエンコードして返す
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(cm.m)
	if err != nil {
		return nil, fmt.Errorf("gob encode error: %w", err)
	}
	return buf.Bytes(), nil
}

func (cm *CacheMap[K, V]) GobDecode(data []byte) error {
	cm.Lock()
	defer cm.Unlock()
	// dataをGobデコードしてcm.mにセットする
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&cm.m); err != nil {
		return fmt.Errorf("gob decode error: %w", err)
	}
	return nil
}
