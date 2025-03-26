package internal

import (
	"fmt"
	"phoenixcache/utils"
	"strings"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
)

// Cache es la estructura que gestiona la caché
type Cache struct {
	store      *ristretto.Cache
	expiration sync.Map
}

type CacheEntry struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	ExpiresIn string `json:"expires_in"`
}

var CacheMutex sync.Mutex

// NewCache crea una nueva instancia de caché
func NewCache(numCounters, maxCost, bufferItems int64) *Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: bufferItems,
	})
	if err != nil {
		panic(fmt.Sprintf("❌ Error al inicializar la caché: %v", err))
	}

	return &Cache{store: cache}
}

//********************************************************************
// Funciones básicas para el uso de ristretto (cache)
//********************************************************************

// Set almacena un valor en la caché con un TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.store.SetWithTTL(key, value, 1, ttl)
	c.expiration.Store(key, time.Now().Add(ttl))
	c.store.Wait()
}

// Get obtiene un valor de la caché si no ha expirado
func (c *Cache) Get(key string) (interface{}, bool) {
	val, found := c.store.Get(key)
	if !found {
		return nil, false
	}

	return val, true
}

func (c *Cache) GetWithExpiry(key string) (interface{}, any, bool) {
	val, found := c.store.Get(key)
	if !found {
		return nil, nil, false
	}

	expTime, exists := c.expiration.Load(key)
	if exists && time.Now().After(expTime.(time.Time)) {
		c.store.Del(key)
		c.expiration.Delete(key)
		return nil, nil, false
	}

	return val, expTime, true
}

// FlushAll borra toda la caché
func (c *Cache) FlushAll() {
	c.store.Clear()
	c.expiration = sync.Map{}
}

// Elimina una Key concreta de la cache
func (c *Cache) RemoveKey(key string) {
	c.store.Del(key)
	c.expiration.Delete(key)
}

func (c *Cache) RemovePatternKey(keyPattern string) []string {

	deletedKeys := []string{}
	// Recorrer la caché y eliminar los que coincidan con el patrón
	c.expiration.Range(func(key, _ interface{}) bool {
		keyStr := key.(string)
		if strings.Contains(keyStr, keyPattern) {
			c.RemoveKey(keyStr)
			deletedKeys = append(deletedKeys, keyStr)
		}
		return true
	})

	return deletedKeys
}

// List devuelve una lista de claves, sus valores truncados y sus expiraciones
func (c *Cache) List() []CacheEntry {
	return c.GetAll(true)
}

// Obtiene la lista de valores (CacheEntry) para realizar la exportación de la cache
func (c *Cache) GetAll(truncateValue bool) []CacheEntry {
	var items []CacheEntry
	c.expiration.Range(func(key, value interface{}) bool {
		val, found := c.store.Get(key.(string))
		if !found {
			return true
		}

		expTime, ok := value.(time.Time)
		if !ok {
			return true
		}

		cacheValue := val.(string)

		if truncateValue {
			cacheValue = utils.TruncateString(cacheValue, 25)
		}

		timeRemaining := time.Until(expTime).String()

		items = append(items, CacheEntry{
			Key:       key.(string),
			Value:     cacheValue,
			ExpiresIn: timeRemaining,
		})
		return true
	})
	return items
}

func (c *Cache) GetDiff() map[string]int64 {
	diff := make(map[string]int64)

	c.expiration.Range(func(key, value interface{}) bool {
		expTime, ok := value.(time.Time)
		if ok {
			diff[key.(string)] = expTime.Unix()
		}
		return true
	})

	return diff
}
