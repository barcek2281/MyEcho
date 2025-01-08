package apiserver

import (
	"net/http"

	"github.com/barcek2281/MyEcho/internal/app/controller"
	"github.com/barcek2281/MyEcho/internal/app/mail"
	"github.com/barcek2281/MyEcho/internal/app/middleware"
	storage "github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type server struct {
	router          *mux.Router
	Logger          *logrus.Logger
	storage         *storage.Storage
	Session         sessions.Store
	controller      *controller.Controller
	controllerPost  *controller.ControllerPost
	controllerAdmin *controller.ControllerAdmin
	middleware      *middleware.Middleware
	Env             Env
}

type ctxKey int8

const (
	sessionName        = "MyEcho"
	ctxKeyUser  ctxKey = iota
)

func newServer(store *storage.Storage, session sessions.Store, logger *logrus.Logger, sender *mail.Sender) *server {
	s := &server{
		router:          mux.NewRouter(),
		Logger:          logger,
		storage:         store,
		Session:         session,
		controller:      controller.NewController(store, session, logger, sender),
		controllerPost:  controller.NewControllerPost(store, session, logger),
		controllerAdmin: controller.NewControllerUser(store, session, logger, sender),
		middleware:      middleware.NewMiddleware(session, store),
	}

	s.ConfigureRouter()
	return s
}

func (s *server) ConfigureRouter() {
	fs := http.FileServer(http.Dir("./static"))
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	s.router.HandleFunc("/", s.controller.MainPage())
	s.router.HandleFunc("/support", s.controller.SupportPage()).Methods("GET")
	s.router.HandleFunc("/support", s.controller.SupportUser()).Methods("POST")

	// мне бы ноормально называть функции, в будущем надо добавить под роутеры :(
	s.router.HandleFunc("/hello", s.controller.HandleHello()).Methods("GET")
	s.router.HandleFunc("/hello", s.controller.HandleHelloPost()).Methods("POST")

	// // Надо будет поменять название функции
	s.router.HandleFunc("/register", s.controller.RegisterUser()).Methods("POST")
	s.router.HandleFunc("/register", s.controller.RegisterPage()).Methods("GET")
	s.router.HandleFunc("/register/verify", s.controller.EmailVerifyPage()).Methods("GET")
	s.router.HandleFunc("/register/verify", s.controller.EmailVerifyUser()).Methods("POST")

	// // я далеко не ушел с названием функций
	s.router.HandleFunc("/login", s.controller.LoginPage()).Methods("GET")
	s.router.HandleFunc("/login", s.controller.LoginUser()).Methods("POST")

	s.router.HandleFunc("/logout", s.controller.LogoutHandler())

	admin := s.router.PathPrefix("/admin").Subrouter()
	admin.PathPrefix("/static/").Handler(http.StripPrefix("/admin/static/", fs))
	admin.Use(s.middleware.AuthenicateAdmin)
	admin.HandleFunc("/", s.controllerAdmin.AdminLoginPage()).Methods("GET")
	admin.HandleFunc("/login", s.controllerAdmin.AdminLogin()).Methods("POST")
	admin.HandleFunc("/register", s.controllerAdmin.AdminRegister()).Methods("POST")
	admin.HandleFunc("/users", s.controllerAdmin.GetAllUsers()).Methods("GET")
	admin.HandleFunc("/updateUserLogin", s.controllerAdmin.UpdateUser()).Methods("POST")
	admin.HandleFunc("/deleteUser", s.controllerAdmin.DeleteUser()).Methods("POST")
	admin.HandleFunc("/findUser", s.controllerAdmin.FindUser()).Methods("POST")
	admin.HandleFunc("/sendMessage", s.controllerAdmin.SendMessageAdmin()).Methods("POST")

	// Лучше его так оставить
	s.router.HandleFunc("/getPost", s.controllerPost.GetPost()).Methods("GET")

	postUrl := s.router.PathPrefix("/post").Subrouter()
	postUrl.Use(s.middleware.AuthenicateUser)
	postUrl.HandleFunc("/createPost", s.controllerPost.CreatePostPage()).Methods("GET")
	postUrl.HandleFunc("/createPost", s.controllerPost.CreatePostReal()).Methods("POST")

}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
