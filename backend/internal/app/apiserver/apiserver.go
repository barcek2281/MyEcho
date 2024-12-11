package apiserver

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
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
}

// --CHAT-GPT
func (s *APIserver) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("handle /hello URL")

		// Подготовим данные для отправки в формате JSON
		response := map[string]string{
			"message": "Hello World!",
		}

		// Установим заголовок Content-Type для ответа
		w.Header().Set("Content-Type", "application/json")
		//io.WriteString(w, "HI")

		// Сериализуем данные в JSON и отправим их клиенту
		if err := json.NewEncoder(w).Encode(response); err != nil {
			// В случае ошибки логируем её и возвращаем HTTP-ошибку
			s.logger.Error("failed to write JSON response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
