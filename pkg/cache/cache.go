/*
 * Copyright 2019 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cache

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	cache *cache.Cache
}

type Item struct {
	Key   string
	Value []byte
}

var ErrNotFound = errors.New("key not found in cache")

func NewCache(defaultExpirationInSeconds int) *Cache {
	return &Cache{cache: cache.New(time.Duration(defaultExpirationInSeconds)*time.Second, time.Duration(defaultExpirationInSeconds)*time.Second)}
}

func (this *Cache) get(key string) (value []byte, err error) {
	temp, found := this.cache.Get(key)
	if !found {
		err = ErrNotFound
	} else {
		var ok bool
		value, ok = temp.([]byte)
		if !ok {
			err = errors.New("unable to interprete cache result")
		}
	}
	return
}

// Add an item to the cache, replacing any existing item.
// If the expiration is 0 the cache's default expiration time is used.
// If it is -1 the item never expires.
func (this *Cache) set(key string, value []byte, expiration time.Duration) {
	this.cache.Set(key, value, expiration)
	return
}

func (this *Cache) Use(key string, expiration time.Duration, getter func() (interface{}, error), result interface{}) (err error) {
	value, err := this.get(key)
	if err == nil {
		err = json.Unmarshal(value, result)
		return
	} else if !errors.Is(err, ErrNotFound) {
		log.Println("WARNING: err in LocalCache::cache.Get()", err)
	}
	temp, err := getter()
	if err != nil {
		return err
	}
	value, err = json.Marshal(temp)
	if err != nil {
		return err
	}
	this.set(key, value, expiration)
	return json.Unmarshal(value, &result)
}

func (this *Cache) UseWithExpirationInResult(key string, getter func() (interface{}, time.Duration, error), result interface{}) (err error) {
	value, err := this.get(key)
	if err == nil {
		err = json.Unmarshal(value, result)
		return
	} else if !errors.Is(err, ErrNotFound) {
		log.Println("WARNING: err in cache.Get()", err)
	}
	temp, expiration, err := getter()
	if err != nil {
		return err
	}
	value, err = json.Marshal(temp)
	if err != nil {
		return err
	}
	this.set(key, value, expiration)
	return json.Unmarshal(value, &result)
}
