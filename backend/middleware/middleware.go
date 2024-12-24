package middleware

import (
	"context"
	"net/http"

	storage "github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/gorilla/sessions"
)

const (
	sessionName = "MyEcho"
	ctxKeyUser  = 1
)

type Middleware struct {
	session sessions.Store
	storage *storage.Storage
}

func NewMiddleware(session sessions.Store, storage *storage.Storage) *Middleware {
	return &Middleware{
		session: session,
		storage: storage,
	}
}
func (m *Middleware) AuthenicateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessions, err := m.session.Get(r, sessionName)
		if err != nil {
			w.Write([]byte("you cant be here"))
			return
		}
		id, ok := sessions.Values["user_id"]
		if !ok {
			w.Write([]byte("you cant be here"))
			return
		}
		u, err := m.storage.User().FindById(id.(int))

		if err != nil {
			w.Write([]byte("you cant be here"))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), 1, u)))
	})
}
