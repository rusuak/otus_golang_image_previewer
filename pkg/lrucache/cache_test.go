package lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Parallel()

	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		t.Parallel()

		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		t.Parallel()

		c := NewCache(5)
		c.Set("aaa", 100)
		c.Set("bbb", 222)

		c.Clear()

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)

		cache := c.(*lruCache)
		require.Zero(t, cache.queue.Len())
		require.Zero(t, len(cache.items))

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
	})

	t.Run("overflow push out", func(t *testing.T) {
		t.Parallel()

		c := NewCache(3)

		c.Set("aaa", 111)
		c.Set("bbb", 222)
		c.Set("ccc", 333)

		// Проверяем элемент напрямую в списке, т.к. Get выносит элемент в начало списка
		lastListItem := c.(*lruCache).queue.Back()
		lastCachedItem := lastListItem.Value.(cacheItem)
		require.Equal(t, "aaa", lastCachedItem.key)
		require.Equal(t, 111, lastCachedItem.value)

		c.Set("ddd", 444)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("last used push out", func(t *testing.T) {
		t.Parallel()

		c := NewCache(3)

		c.Set("aaa", 111) // [aaa]
		c.Set("bbb", 222) // [bbb,aaa]
		c.Set("ccc", 333) // [ccc,bbb,aaa]

		_, ok := c.Get("ccc") // [ccc,bbb,aaa]
		require.True(t, ok)
		_, ok = c.Get("bbb") // [bbb,ccc,aaa]
		require.True(t, ok)
		_, ok = c.Get("aaa") // [aaa,bbb,ccc]
		require.True(t, ok)
		c.Get("ccc")

		c.Set("ddd", 444)

		val, ok := c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})
}
