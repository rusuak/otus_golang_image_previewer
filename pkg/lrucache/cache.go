package lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	newValue := cacheItem{key, value}

	if _, ok := cache.items[key]; ok {
		listItem := cache.items[key]
		listItem.Value = newValue
		cache.queue.MoveToFront(listItem)

		return true
	}

	newItem := cache.queue.PushFront(newValue)
	cache.items[key] = newItem

	if cache.queue.Len() > cache.capacity {
		cache.removeLastItem()
	}

	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	if _, ok := cache.items[key]; !ok {
		return nil, false
	}

	listItem := cache.items[key]
	cache.queue.MoveToFront(listItem)

	cachedItem := listItem.Value.(cacheItem)
	return cachedItem.value, true
}

func (cache *lruCache) Clear() {
	c := NewCache(cache.capacity)
	emptyCachePtr := c.(*lruCache)
	*cache = *emptyCachePtr
}

func (cache *lruCache) removeLastItem() {
	lastItem := cache.queue.Back()
	cache.queue.Remove(lastItem)

	cachedItem := lastItem.Value.(cacheItem)
	delete(cache.items, cachedItem.key)
}
