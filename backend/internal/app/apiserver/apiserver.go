// package apiserver

// import (
// 	"net/http"

// 	"github.com/barcek2281/MyEcho/internal/app/storage"

// 	"github.com/gorilla/mux"
// 	"github.com/sirupsen/logrus"
// )

// type APIserver struct {
// 	config     *Config
// 	Logger     *logrus.Logger
// 	router     *mux.Router
// 	controller *Controller
// 	storage    *storage.Storage
// }

// func NewAPIserver(config *Config) *APIserver {
// 	return &APIserver{
// 		config:     config,
// 		Logger:     logrus.New(),
// 		router:     mux.NewRouter(),
// 		controller: NewController(),
// 	}
// }

// func (s *APIserver) Start() error {
// 	if err := s.ConfigureLogger(); err != nil {
// 		return err
// 	}

// 	s.ConfigureRouter()

// 	if err := s.ConfigureStorage(); err != nil {
// 		s.Logger.Errorf("failed to configure storage: %v", err)
// 		return err
// 	}

// 	s.Logger.Info("Starting API server: http://localhost", s.config.BinAddr)

// 	return http.ListenAndServe(s.config.BinAddr, s.router)
// }

// func (s *APIserver) ConfigureLogger() error {
// 	// Обычная конфигурация Логгера
// 	level, err := logrus.ParseLevel(s.config.LogLevel)
// 	if err != nil {
// 		return err
// 	}
// 	s.Logger.SetLevel(level)
// 	return nil
// }

// func (s *APIserver) ConfigureRouter() {
// 	s.router.HandleFunc("/", s.controller.MainPage(s))
// 	// мне бы ноормально называть функции, в будущем надо добавить под роутеры :(
// 	s.router.HandleFunc("/hello", s.controller.handleHello(s)).Methods("GET")
// 	s.router.HandleFunc("/hello", s.controller.handleHelloPost(s)).Methods("POST")

// 	// Надо будет поменять название функции
// 	s.router.HandleFunc("/register", s.controller.registerUser(s)).Methods("POST")
// 	s.router.HandleFunc("/register", s.controller.registerPage(s)).Methods("GET")

// 	s.router.HandleFunc("/users", s.controller.getAllUsers(s)).Methods("GET")
// 	s.router.HandleFunc("/updateUserLogin", s.controller.UpdateUser(s)).Methods("POST")
// 	s.router.HandleFunc("/deleteUser", s.controller.DeleteUser(s)).Methods("POST")
// 	s.router.HandleFunc("/findUser", s.controller.FindUser(s)).Methods("POST")
// }

// func (s *APIserver) ConfigureStorage() error {
// 	st := storage.New(s.config.DataBaseURL)
// 	if err := st.Open(); err != nil { // Ping db
// 		return err
// 	}
// 	s.storage = st
// 	return nil
// }

package apiserver

import (
	"net/http"
	"os"

	"github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/barcek2281/MyEcho/mail"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func Start(config *Config) error {
	store := storage.New(config.DataBaseURL)
	if err := store.Open(); err != nil { // Ping db
		return err
	}
	session := sessions.NewCookieStore([]byte(config.CookieKey))

	logger := logrus.New()
	logger.SetFormatter(&prefixed.TextFormatter{
		DisableColors: true,
		TimestampFormat : "2006-01-02 15:04:05",
		FullTimestamp:true,
		ForceFormatting: true,
	})
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	} else {
		logger.Out = f
	}

	logger.SetLevel(level)

	sender := mail.NewSender(config.EmailTo, config.EmailToPassword)

	s := newServer(store, session, logger, sender)
	return http.ListenAndServe(config.BinAddr, s)
}