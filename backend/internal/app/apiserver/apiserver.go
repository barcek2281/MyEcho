package apiserver

import (
	"github.com/barcek2281/MyEcho/internal/app/storage"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIserver struct {
	config     *Config
	Logger     *logrus.Logger
	router     *mux.Router
	controller *Controller
	storage    *storage.Storage
}

func NewAPIserver(config *Config) *APIserver {
	return &APIserver{
		config:     config,
		Logger:     logrus.New(),
		router:     mux.NewRouter(),
		controller: NewController(),
	}
}

func (s *APIserver) Start() error {
	if err := s.ConfigureLogger(); err != nil {
		return err
	}

	s.ConfigureRouter()

	if err := s.ConfigureStorage(); err != nil {
		s.Logger.Errorf("failed to configure storage: %v", err)
		return err
	}

	s.Logger.Info("Starting API server: http://localhost", s.config.BinAddr)

	return http.ListenAndServe(s.config.BinAddr, s.router)
}

func (s *APIserver) ConfigureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.Logger.SetLevel(level)
	return nil
}

func (s *APIserver) ConfigureRouter() {
	s.router.HandleFunc("/", s.controller.MainPage(s))
	s.router.HandleFunc("/hello", s.controller.handleHello(s))
	s.router.HandleFunc("/register", s.controller.register(s))
}

func (s *APIserver) ConfigureStorage() error {
	st := storage.New(s.config.DataBaseURL)
	if err := st.Open(); err != nil {
		return err
	}
	s.storage = st
	return nil
}
