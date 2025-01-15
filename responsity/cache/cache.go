package cache

import (
	"MiniPrograms/responsity/model"
	"sync"
)

// Cache 定义一个线程安全的缓存
type Cache struct {
	cache sync.Map
}

// NewCache 创建缓存实例
func NewCache() *Cache {
	return &Cache{}
}

// Load 从缓存中获取数据，返回指针和存在标记
func (c *Cache) Load(key string) (*model.MiniPrograms, bool) {
	if value, ok := c.cache.Load(key); ok {
		// 类型断言，并直接返回,非常清晰而优雅的设计
		if resp, ok := value.(*model.MiniPrograms); ok {
			return resp, true
		}
	}
	return nil, false
}

// Store 将数据存入缓存
func (c *Cache) Store(key string, resp *model.MiniPrograms) {
	c.cache.Store(key, resp)
}
