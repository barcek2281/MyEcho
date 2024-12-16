package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/barcek2281/MyEcho/internal/app/model"
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

		response := map[string]string{
			"status":  "OK",
			"message": "Hello World!",
		}
		s.Logger.Warn("handle /hello")

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

func (ctrl *Controller) registerPage(s *APIserver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/app/templates/register.html")
		if err != nil {
			s.Logger.Error(err)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			s.Logger.Error(err)
			return
		}
		s.Logger.Info("handle /register GET")
	}
}

func (ctrl *Controller) registerUser(s *APIserver) http.HandlerFunc {
	type Request struct {
		Login    string `json:"login"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.Logger.Error(err)
			return
		}

		u := model.User{
			Email:    req.Email,
			Password: req.Password,
			Login:    req.Login,
		}

		if err := s.storage.User().Create(&u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.Logger.Error(err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		s.Logger.Info("handle /register POST")
	}
}
