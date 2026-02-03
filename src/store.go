package main

import (
	"sync"
	"time"
)

const seedContent = `<h2>Welcome to the editor</h2>
<p>This content is stored in memory.</p>
<blockquote><p>Edit this and click Save.</p></blockquote>`

type contentStore struct {
	mu    sync.RWMutex
	items map[string]contentEntry
}

type contentEntry struct {
	content  string
	lastSeen time.Time
}

func newContentStore() *contentStore {
	return &contentStore{
		items: make(map[string]contentEntry),
	}
}

func (store *contentStore) loadContent(clientID string) (string, error) {
	store.mu.RLock()
	entry, ok := store.items[clientID]
	store.mu.RUnlock()
	if ok {
		store.touch(clientID)
		return entry.content, nil
	}
	return seedContent, nil
}

func (store *contentStore) saveContent(clientID string, html string) {
	store.mu.Lock()
	store.items[clientID] = contentEntry{
		content:  html,
		lastSeen: time.Now(),
	}
	store.mu.Unlock()
}

func (store *contentStore) touch(clientID string) {
	store.mu.Lock()
	entry, ok := store.items[clientID]
	if ok {
		entry.lastSeen = time.Now()
		store.items[clientID] = entry
	}
	store.mu.Unlock()
}

func (store *contentStore) pruneExpired(ttl time.Duration) int {
	store.mu.Lock()
	defer store.mu.Unlock()

	if ttl <= 0 {
		return 0
	}

	cutoff := time.Now().Add(-ttl)
	pruned := 0
	for clientID, entry := range store.items {
		if entry.lastSeen.Before(cutoff) {
			delete(store.items, clientID)
			pruned++
		}
	}

	return pruned
}
