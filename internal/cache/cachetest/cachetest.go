// Package cachetest provides an in-memory cache.Interface for tests.
package cachetest

import (
	"context"
	"errors"
	"sync"

	"github.com/Template7/backend/internal/cache"
)

// ErrMiss mimics redis.Nil for a missing key.
var ErrMiss = errors.New("cache miss")

type Memory struct {
	mu    sync.Mutex
	codes map[string]string
}

var _ cache.Interface = (*Memory)(nil)

func New() *Memory {
	return &Memory{codes: map[string]string{}}
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
