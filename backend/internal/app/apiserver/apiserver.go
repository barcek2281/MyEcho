package apiserver

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIserver struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

func NewAPIserver(config *Config) *APIserver {
	return &APIserver{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIserver) Start() error {
	if err := s.ConfigureLogger(); err != nil {
		return err
	}
	s.ConfigureRouter()

	s.logger.Info("Starting API server")

	return http.ListenAndServe(s.config.BinAddr, s.router)
}

func (s *APIserver) ConfigureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIserver) ConfigureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
	s.router.HandleFunc("/post", s.post())
}

// --CHAT-GPT
func (s *APIserver) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Установим заголовок Content-Type для ответа
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"status":  "OK",
			"message": "Hello World!",
		}

		switch r.Method {
		case "GET":
			s.logger.Info("handle /hello GET")

			name := r.URL.Query()
			prettyJson, err := json.MarshalIndent(name, "", "  ")
			if err != nil {
				s.logger.Error(err)
			}
			fmt.Println(string(prettyJson))
			if err := json.NewEncoder(w).Encode(response); err != nil {
				s.logger.Error(err)
				http.Error(w, "lol", http.StatusInternalServerError)
			}
		default:
			s.logger.Info("Unhandled Unknown method /hello")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (s *APIserver) post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			s.logger.Warn("Unhandled method " + r.Method)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			s.logger.Error(err)
		}

		prettyJson, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			s.logger.Error(err)
		}
		fmt.Println(string(prettyJson))
	}
}
