package session

import (
	"encoding/gob"
	"main/config"
	"math"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

type Session struct {
	Store *sessions.FilesystemStore
}

func NewSession(conf config.ConfigManager) Session {
	s := sessions.NewFilesystemStore(os.TempDir(), []byte(conf.GetSessionKey()))
	s.MaxLength(math.MaxInt64)
	gob.Register(map[string]interface{}{})
	return Session{
		Store: s,
	}
}

func (s *Session) Get(r *http.Request, sessionName string) (*sessions.Session, error) {
	return s.Store.Get(r, sessionName)
}
