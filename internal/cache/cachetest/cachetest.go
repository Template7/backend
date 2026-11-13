// Package cachetest provides an in-memory cache.Interface for tests.
package cachetest

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Template7/backend/internal/cache"
)

// ErrMiss mimics redis.Nil for a missing key.
var ErrMiss = errors.New("cache miss")

type Memory struct {
	mu       sync.Mutex
	codes    map[string]string
	failures map[string]counter
}

type counter struct {
	n       int64
	expires time.Time
}

var _ cache.Interface = (*Memory)(nil)

func New() *Memory {
	return &Memory{codes: map[string]string{}, failures: map[string]counter{}}
}

func (m *Memory) SetUserActivationCode(_ context.Context, userId string, code string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.codes[userId] = code
	return nil
}

func (m *Memory) GetUserActivationCode(_ context.Context, userId string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	code, ok := m.codes[userId]
	if !ok {
		return "", ErrMiss
	}
	return code, nil
}

func (m *Memory) DeleteUserActivationCode(_ context.Context, userId string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.codes, userId)
	return nil
}

// Has reports whether an activation code is stored for userId.
func (m *Memory) Has(userId string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.codes[userId]
	return ok
}

func (m *Memory) IncrLoginFailures(_ context.Context, key string, window time.Duration) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c := m.failures[key]
	if c.n == 0 || time.Now().After(c.expires) {
		c = counter{expires: time.Now().Add(window)}
	}
	c.n++
	m.failures[key] = c
	return c.n, nil
}

func (m *Memory) GetLoginFailures(_ context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c := m.failures[key]
	if time.Now().After(c.expires) {
		return 0, nil
	}
	return c.n, nil
}

func (m *Memory) ResetLoginFailures(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.failures, key)
	return nil
}
