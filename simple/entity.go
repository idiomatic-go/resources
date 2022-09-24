package util

import (
	"hash/fnv"
	"sync"
	"sync/atomic"
)

type GetVersion func(any) string

type VersionedEntity struct {
	hash   uint32
	entity any
	mu     sync.RWMutex
	fn     GetVersion
}

func CreateVersionedEntity(entity any, fn GetVersion) *VersionedEntity {
	get := fn
	if get == nil {
		get = func(entity any) string { return "" }
	}
	// Dereference as all helper functions can use the underlying struct
	if IsPointer(entity) {
		entity = Copy(entity)
	}
	ve := &VersionedEntity{hash: 0, entity: nil, fn: get}
	ve.Set(entity)
	return ve
}

func (v *VersionedEntity) IsNewVersion(version string) bool {
	if v.entity == nil {
		return false
	}
	return hash(version) != atomic.LoadUint32(&v.hash)
}

// Set - set the entity
func (v *VersionedEntity) Set(entity any) {
	v.mu.Lock()
	v.entity = entity
	if v.entity == nil {
		atomic.StoreUint32(&v.hash, 0)
	} else {
		atomic.StoreUint32(&v.hash, hash(v.fn(v.entity)))
	}
	v.mu.Unlock()
}

// Get - get a copy of the entity
func (v *VersionedEntity) Get() any {
	if v.entity == nil {
		return nil
	}
	v.mu.RLock()
	c := Copy(v.entity)
	v.mu.RUnlock()
	return c
}

// GetVersion - get the version string
func (v *VersionedEntity) GetVersion() string {
	if v.entity == nil {
		return ""
	}
	v.mu.RLock()
	vers := v.fn(v.entity)
	v.mu.RUnlock()
	return vers
}

func hash(s string) uint32 {
	if s == "" {
		return 0
	}
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
