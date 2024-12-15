package apiserver

import (
	"encoding/json"
	"fmt"
	_ "github.com/sirupsen/logrus"
	"net/http"
	"text/template"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) MainPage(s *APIserver) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//data := "Go Template"
		tmpl, err := template.ParseFiles("internal/app/templates/index.html")
		if err != nil {
			s.Logger.Error(err)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			s.Logger.Error(err)
			return
		}
	}
}

func (ctrl *Controller) handleHello(s *APIserver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			s.Logger.Warn("Handel /hello: Method not allowed")
			return
		}

		response := map[string]string{
			"status":  "OK",
			"message": "Hello World!",
		}
		s.Logger.Info("handle /hello GET")

		// Print query
		prettyJson, err := json.MarshalIndent(r.URL.Query(), "", "  ")
		if err != nil {
			s.Logger.Error(err)
		}
		fmt.Println(string(prettyJson))

		// receive accept
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			s.Logger.Error(err)
			http.Error(w, "cannot write json file", http.StatusInternalServerError)
		}
	}
}

func (ctrl *Controller) register(s *APIserver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			s.Logger.Warn("Unhandled method " + r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			s.Logger.Error(err)
		}

		prettyJson, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			s.Logger.Error(err)
		}
		fmt.Println(string(prettyJson))
	}
}
