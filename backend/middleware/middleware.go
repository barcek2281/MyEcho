package middleware

import (
	"context"
	"errors"
	"net/http"

	storage "github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/barcek2281/MyEcho/pkg/utils"
	"github.com/gorilla/sessions"
)

const (
	sessionName = "MyEcho"
	sessionAdmin = "IsAdmin"
	ctxKeyUser  = 1
)

var (
	errYouCantBeHere = errors.New("You cant be here")
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
		// TODO: refactor errors
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

func (m *Middleware) AuthenicateAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedPaths := map[string]struct{}{
			"/admin/":         {},
			"/admin/login":    {},
			"/admin/register": {},
		}
		if _, ok := allowedPaths[r.URL.Path]; ok {
			// Пропускаем без проверки
			next.ServeHTTP(w, r)
			return
		}
		session, err := m.session.Get(r, sessionAdmin)
		id, ok := session.Values["admin_id"]
		if !ok {
			utils.Error(w, r, http.StatusBadGateway, errYouCantBeHere)
			return
		}
		a, err := m.storage.Admin().FindById(id.(int))

		if err != nil {
			utils.Error(w, r, http.StatusBadGateway, errYouCantBeHere)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), 1, a)))

	})
}
