package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
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
		keys1 := []Key{"aaa", "bbb", "ccc", "ddd", "eee"}
		keys2 := []Key{"foo", "bar", "zoo", "van"}

		for _, oldestKey := range keys1 {
			c := NewCache(len(keys1))
			for _, k := range keys1 {
				c.Set(k, 1)
			}

			c.Get(oldestKey)
			for _, k := range keys2 { // purge all keys except the one
				c.Set(k, 1)
			}
			_, ok := c.Get(oldestKey)
			require.True(t, ok)
		}
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache(3)
		for i, k := range []string{"aaa", "bbb", "ccc"} {
			c.Set(Key(k), i)
		}

		c.Clear()
		for _, k := range []string{"aaa", "bbb", "ccc"} {
			_, ok := c.Get(Key(k))
			require.False(t, ok)
		}
	})

	t.Run("zero sized cache", func(t *testing.T) {
		c := NewCache(0)

		wasInCache := c.Set("aaa", 1)
		require.False(t, wasInCache)

		_, ok := c.Get("aaa")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := range 1_000_000 {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for range 1_000_000 {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
