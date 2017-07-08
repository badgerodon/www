package rbsa

import (
	"time"
)

type (
	Cache struct {
		lock   chan bool
		lookup map[string]CacheEntry
		size   int
	}
	CacheEntry struct {
		access time.Time
		data   interface{}
	}
)

var LEAST time.Time

func NewCache(size int) *Cache {
	c := &Cache{make(chan bool, 1), make(map[string]CacheEntry), size}

	c.lock <- true

	return c
}
func (this *Cache) Get(key string, fill func() (interface{}, error)) (interface{}, error) {
	<-this.lock
	entry, ok := this.lookup[key]
	if ok {
		entry.access = time.Now()
		this.lock <- true
		return entry.data, nil
	}
	this.lock <- true

	v, err := fill()
	if err != nil {
		return nil, err
	}
	<-this.lock
	this.lookup[key] = CacheEntry{time.Now(), v}
	if len(this.lookup) > this.size {
		least := ""
		var last time.Time
		for k, e := range this.lookup {
			if last == LEAST || e.access.Before(last) {
				last = e.access
				least = k
			}
		}

		if least != "" {
			delete(this.lookup, least)
		}
	}
	this.lock <- true

	return v, nil
}
