package apiserver

import (
	"encoding/json"

	"net/http"

	"github.com/barcek2281/MyEcho/internal/app/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type server struct {
	router     *mux.Router
	Logger     *logrus.Logger
	storage    *storage.Storage
	Session    sessions.Store
	controller *Controller
}

func newServer(store *storage.Storage, session sessions.Store) *server {
	s := &server{
		router:  mux.NewRouter(),
		Logger:  logrus.New(),
		storage: store,
		Session: session,

	}
	s.ConfigureRouter()

	return s
}

func (s *server) ConfigureRouter() {
	s.router.HandleFunc("/", controller.MainPage(s))

	// мне бы ноормально называть функции, в будущем надо добавить под роутеры :(
	s.router.HandleFunc("/hello", controller.handleHello(s)).Methods("GET")
	s.router.HandleFunc("/hello", controller.handleHelloPost(s)).Methods("POST")

	// Надо будет поменять название функции
	s.router.HandleFunc("/register", controller.registerUser(s)).Methods("POST")
	s.router.HandleFunc("/register", controller.registerPage(s)).Methods("GET")

	s.router.HandleFunc("/users", controller.getAllUsers(s)).Methods("GET")
	s.router.HandleFunc("/updateUserLogin", controller.UpdateUser(s)).Methods("POST")
	s.router.HandleFunc("/deleteUser", controller.DeleteUser(s)).Methods("POST")
	s.router.HandleFunc("/findUser", controller.FindUser(s)).Methods("POST")
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.Respond(w, r, code, map[string]string{"error": err.Error()})

}
func (s *server) Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
