package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type Session struct {
	UserID    uint
	Phone     string
	Role      string
	ExpiresAt time.Time
}

type AuthManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewAuthManager() *AuthManager {
	return &AuthManager{
		sessions: make(map[string]*Session),
	}
}

// Generate token session
func (a *AuthManager) CreateSession(userID uint, phone, role string) string {
	a.mu.Lock()
	defer a.mu.Unlock()

	token := generateToken()
	a.sessions[token] = &Session{
		UserID:    userID,
		Phone:     phone,
		Role:      role,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	return token
}

func (a *AuthManager) GetSession(token string) (*Session, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	session, exists := a.sessions[token]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil, false
	}
	return session, true
}

func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
