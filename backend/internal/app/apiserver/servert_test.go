package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barcek2281/MyEcho/internal/app/controller"
	storage "github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	config := NewConfig()
	config.LogFilePath = "../../../log/info.log"
	store := storage.New(config.DataBaseURL)
	if err := store.Open(); err != nil { // Ping db
		assert.NotNil(t, err)
	}
	session := sessions.NewCookieStore([]byte(config.CookieKey))

	ctrl := controller.NewController(store, session, logrus.New(), nil)

	ctrl.HandleHello()(w, req)

	res := w.Result()
	if res.StatusCode != 200 {
		t.Error("doesnt work")
	}
}

func ReturnController(logfilePath string) *controller.Controller {
	config := NewConfig()
	config.LogFilePath = logfilePath
	store := storage.New(config.DataBaseURL)
	if err := store.Open(); err != nil { // Ping db
		return nil
	}
	session := sessions.NewCookieStore([]byte(config.CookieKey))

	ctrl := controller.NewController(store, session, logrus.New(), nil)
	return ctrl
}
