package main

import "sync"

const seedContent = `<h2>Welcome to the editor</h2>
<p>This content is stored in memory.</p>
<blockquote><p>Edit this and click Save.</p></blockquote>`

type contentStore struct {
	mu    sync.RWMutex
	items map[string]string
}

func newContentStore() *contentStore {
	return &contentStore{
		items: make(map[string]string),
	}
}

func loadContent(store *contentStore, clientID string) (string, error) {
	store.mu.RLock()
	content, ok := store.items[clientID]
	store.mu.RUnlock()
	if ok {
		return content, nil
	}
	return seedContent, nil
}

func saveContent(store *contentStore, clientID string, html string) {
	store.mu.Lock()
	store.items[clientID] = html
	store.mu.Unlock()
}
