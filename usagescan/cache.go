package usagescan

import (
	"os"
	"strings"
	"sync"
)

var GlobalCache = &DiskCache{
	entries: make(map[string]*CacheEntry),
}

type DiskCache struct {
	sync.RWMutex
	entries map[string]*CacheEntry
}

type CacheEntry struct {
	Path      string
	Size      int64
	Done      bool
	mu        sync.Mutex
	subs      map[uint64]func(int64)
	nextSubID uint64
	doneCh    chan struct{}
}

func (c *DiskCache) GetEntry(path string) *CacheEntry {
	c.RLock()
	defer c.RUnlock()
	return c.entries[path]
}

func (c *DiskCache) GetOrCreateEntry(path string) (*CacheEntry, bool) {
	c.Lock()
	defer c.Unlock()
	entry, exists := c.entries[path]
	if !exists {
		entry = &CacheEntry{
			Path:   path,
			subs:   make(map[uint64]func(int64)),
			doneCh: make(chan struct{}),
		}
		c.entries[path] = entry
	}
	return entry, exists
}

func (c *DiskCache) Invalidate(path string) {
	c.Lock()
	defer c.Unlock()

	separator := string(os.PathSeparator)
	prefix := path
	if !strings.HasSuffix(path, separator) {
		prefix = path + separator
	}

	for key := range c.entries {
		if key == path || strings.HasPrefix(key, prefix) {
			delete(c.entries, key)
		}
	}
}

func (e *CacheEntry) Subscribe(onProgress func(int64)) (unsubscribe func()) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.Done {
		onProgress(e.Size)
		return func() {}
	}

	id := e.nextSubID
	e.nextSubID++
	e.subs[id] = onProgress

	onProgress(e.Size)

	return func() {
		e.mu.Lock()
		delete(e.subs, id)
		e.mu.Unlock()
	}
}

func (e *CacheEntry) UpdateSize(size int64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Size = size
	for _, sub := range e.subs {
		sub(size)
	}
}

func (e *CacheEntry) MarkDone() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Done = true
	for _, sub := range e.subs {
		sub(e.Size)
	}
	e.subs = nil
	close(e.doneCh)
}

func (e *CacheEntry) Wait() {
	<-e.doneCh
}