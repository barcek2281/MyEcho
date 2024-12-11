package apiserver

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIserver_HandleHello(t *testing.T) {
	s := NewAPIserver(NewConfig())
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	s.handleHello().ServeHTTP(rec, req)
	assert.Equal(t, rec.Body.String(), "Hello World!")
}
