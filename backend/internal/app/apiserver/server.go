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
	Env        Env
}

func newServer(store *storage.Storage, session sessions.Store, logger *logrus.Logger, env *Env) *server {
	s := &server{
		router:  mux.NewRouter(),
		Logger:  logger,
		storage: store,
		Session: session,
		Env:     *env,
	}
	
	s.ConfigureRouter()
	return s
}

func (s *server) ConfigureRouter() {
	fs := http.FileServer(http.Dir("./static"))
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	s.router.HandleFunc("/", controller.MainPage(s))
	s.router.HandleFunc("/support", controller.SupportPage(s)).Methods("GET")
	s.router.HandleFunc("/support", controller.SupportUser(s)).Methods("POST")

	// мне бы ноормально называть функции, в будущем надо добавить под роутеры :(
	s.router.HandleFunc("/hello", controller.handleHello(s)).Methods("GET")
	s.router.HandleFunc("/hello", controller.handleHelloPost(s)).Methods("POST")

	// Надо будет поменять название функции
	s.router.HandleFunc("/register", controller.registerUser(s)).Methods("POST")
	s.router.HandleFunc("/register", controller.registerPage(s)).Methods("GET")

	// я далеко не ушел с названием функций
	s.router.HandleFunc("/login", controller.loginPage(s)).Methods("GET")
	s.router.HandleFunc("/login", controller.loginUser(s)).Methods("POST")

	s.router.HandleFunc("/logout", controller.LogoutHandler(s))

	// TODO: разделить для админа эти ссылка
	s.router.HandleFunc("/users", controllerUser.getAllUsers(s)).Methods("GET")
	s.router.HandleFunc("/updateUserLogin", controllerUser.UpdateUser(s)).Methods("POST")
	s.router.HandleFunc("/deleteUser", controllerUser.DeleteUser(s)).Methods("POST")
	s.router.HandleFunc("/findUser", controllerUser.FindUser(s)).Methods("POST")

	// TODO: отдельно добавить ссылку для постов
	s.router.HandleFunc("/createPost", controllerPost.CreatePost(s)).Methods("GET")
	s.router.HandleFunc("/createPost", controllerPost.CreatePostReal(s)).Methods("POST")
	s.router.HandleFunc("/getPost", controllerPost.GetPost(s)).Methods("GET")
}


func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.Respond(w, r, code, map[string]string{"error": err.Error()})

}
func (s *server) Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
