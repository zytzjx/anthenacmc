package cmcupdate

import (
	"errors"
	"fmt"
	"sync"

	dmc "github.com/zytzjx/anthenacmc/datacentre"
)

// SyncDownLoadList downlist
type SyncDownLoadList struct {
	mp    map[string][]string
	mutex *sync.Mutex
}

// NewSyncDownLoadList new SyncDownLoadList
func NewSyncDownLoadList() *SyncDownLoadList {
	sdl := SyncDownLoadList{}
	sdl.mp = make(map[string][]string)
	sdl.mutex = &sync.Mutex{}
	return &sdl
}

// Get  values by key
func (sdll *SyncDownLoadList) Get(key string) ([]string, error) {
	sdll.mutex.Lock()
	defer sdll.mutex.Unlock()
	ss, ok := sdll.mp[key]

	if !ok {
		return ss, errors.New("not exist")
	}

	return ss, nil
}

// Set key values
func (sdll *SyncDownLoadList) Set(key string, v []string) {
	sdll.mutex.Lock()
	defer sdll.mutex.Unlock()
	sdll.mp[key] = v
}

// SetItem key append v
func (sdll *SyncDownLoadList) SetItem(key string, v string) {
	sdll.mutex.Lock()
	defer sdll.mutex.Unlock()
	if _, ok := sdll.mp[key]; !ok {
		sdll.mp[key] = []string{}
	}
	sdll.mp[key] = append(sdll.mp[key], v)
}

// Display show print
func (sdll *SyncDownLoadList) Display() {
	sdll.mutex.Lock()
	defer sdll.mutex.Unlock()

	for k, v := range sdll.mp {
		fmt.Println(k, "=", v)
	}
}

// SaveRedis save to redis
func (sdll *SyncDownLoadList) SaveRedis() {
	sdll.mutex.Lock()
	defer sdll.mutex.Unlock()

	for k, v := range sdll.mp {
		for _, vv := range v {
			dmc.AddSet(k, vv)
		}
	}
	dmc.Set("hydradownload.status", "complete", 0)
}

// RemoveRedis remove
func RemoveRedis(key string) {
	dmc.Del(key)
}