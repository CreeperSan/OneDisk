package storage

import (
	errcode "OneDisk/def/err_code"
	defstorage "OneDisk/def/storage"
	"OneDisk/lib/log"
	"OneDisk/module/database"
	"encoding/json"
	"go.uber.org/zap"
	"sync"
)

type CacheStorage struct {
	cache map[int64]PlatformInterface
	lock  sync.RWMutex
}

var instance *CacheStorage
var once sync.Once

func Cache() *CacheStorage {
	once.Do(func() {
		instance = &CacheStorage{
			cache: make(map[int64]PlatformInterface),
		}
	})
	return instance
}

func (c *CacheStorage) SetStorage(storageID int64, storage PlatformInterface) {
	c.lock.Lock()
	c.cache[storageID] = storage
	defer c.lock.Unlock()
}

func (c *CacheStorage) GetStorage(storageID int64) PlatformInterface {
	c.lock.RLock()
	v, exist := c.cache[storageID]
	defer c.lock.RUnlock()
	// 1、缓存存在则返回
	if exist {
		return v
	}
	// 2、缓存不存在，则读取数据库
	tmpStorage, tmpResult := database.StorageFind(storageID)
	if tmpResult.Code != errcode.OK {
		log.Error(tag, "Error occurred while finding storage in Cache.GetStorage()", zap.Error(tmpResult.Error))
	}
	// 2.1、数据库没有找到缓存并返回 nil
	if tmpStorage == nil {
		c.lock.RLock()
		c.cache[storageID] = nil
		defer c.lock.RUnlock()
		return nil
	}
	// 2.2、数据库找到缓存并返回
	if tmpStorage.Type == database.ValueStorageTypePath {
		var configLocalPath defstorage.ConfigLocalPath
		err := json.Unmarshal([]byte(tmpStorage.Config), &configLocalPath)
		if err != nil {
			log.Error(tag, "Error occurred while unmarshal storage config in Cache.GetStorage()", zap.Error(err))
			c.lock.RLock()
			c.cache[storageID] = nil
			defer c.lock.RUnlock()
			return nil
		}
		// 3.2、实例化
		prefabStorage := PlatformInterfaceLocal{
			Root:             configLocalPath.Path,
			DownloadUrlCache: make(map[string]string),
			UploadUrlCache:   make(map[string]string),
		}
		c.lock.RLock()
		c.cache[storageID] = &prefabStorage
		defer c.lock.RUnlock()
		return &prefabStorage
	} else {
		log.Error(tag, "Storage type not supported in Cache.GetStorage()", zap.Int("storageType", tmpStorage.Type))
		c.lock.RLock()
		c.cache[storageID] = nil
		defer c.lock.RUnlock()
		return nil
	}
}
