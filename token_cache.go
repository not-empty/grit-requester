package gritrequester

import "sync"

type TokenCache struct {
	mu     sync.RWMutex
	tokens map[string]string
}

func NewTokenCache() *TokenCache {
	return &TokenCache{tokens: make(map[string]string)}
}

func (c *TokenCache) Get(service string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	t, ok := c.tokens[service]
	return t, ok
}

func (c *TokenCache) Set(service, token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tokens[service] = token
}

func (c *TokenCache) Delete(service string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.tokens, service)
}
