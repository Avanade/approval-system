package session

import (
	"encoding/gob"
	"fmt"
	"main/config"
	"math"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

type Session struct {
	Store *sessions.FilesystemStore
}

func NewSession(conf config.ConfigManager) ConnectSession {
	s := sessions.NewFilesystemStore(os.TempDir(), []byte(conf.GetSessionKey()))
	s.MaxLength(math.MaxInt64)
	gob.Register(map[string]interface{}{})
	return &Session{Store: s}
}

func (s *Session) ConnectSession(w *http.ResponseWriter, r *http.Request, session string) SessionManager {

	return &sessionManager{
		w:           w,
		r:           r,
		values:      make(map[string]interface{}),
		sessionName: session,
		store:       s.Store,
	}
}

// Destroy implements SessionManager.
func (s *sessionManager) Destroy() error {
	session, err := s.store.Get(s.r, s.sessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(s.r, *s.w)
}

// Get implements SessionManager.
func (s *sessionManager) Get(key string) (interface{}, error) {
	session, err := s.store.Get(s.r, s.sessionName)
	if err != nil {
		return nil, err
	}

	if _, ok := session.Values[key]; !ok {
		return nil, fmt.Errorf("no key found in session")
	}
	return session.Values[key], nil
}

// Save implements SessionManager.
func (s *sessionManager) Save(maxAge int) error {
	session, err := s.store.Get(s.r, s.sessionName)
	if err != nil {
		return err
	}
	for k, v := range s.values {
		session.Values[k] = v
	}

	session.Options.MaxAge = maxAge

	return session.Save(s.r, *s.w)
}

// Set implements SessionManager.
func (s *sessionManager) Set(key string, value interface{}) {
	s.values[key] = value
}
