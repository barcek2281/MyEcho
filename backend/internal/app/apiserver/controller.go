package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/barcek2281/MyEcho/internal/app/model"
	_ "github.com/sirupsen/logrus"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) MainPage(s *APIserver) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/index.html")

		if err != nil {
			s.Logger.Error(err)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			s.Logger.Error(err)
			return
		}
		s.Logger.Info("handle MainPage GET")
	}
}

func (ctrl *Controller) handleHello(s *APIserver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response := map[string]string{
			"status":  "OK",
			"message": "Hello World!",
		}
		s.Logger.Info("handle /hello")

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

func (ctrl *Controller) handleHelloPost(s *APIserver) http.HandlerFunc {
	type Request struct {
		Msg string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "not ok", "msg": "we got a cringe message"})
			return
		}

		if req.Msg == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "not ok", "msg": "we got a empty message "})
			return
		}

		fmt.Println(req)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "msg": "we got the message"})

		s.Logger.Info("handle /hello POST")

	}

}

func (ctrl *Controller) registerPage(s *APIserver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/register.html")
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

func (ctrl *Controller) getAllUsers(s *APIserver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := s.storage.User().GetAll(20)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.Logger.Error(err)
			return
		}
		w.Header().Set("Content-Type", "http")
		tmpl, err := template.ParseFiles("./templates/users.html")
		if err != nil {
			s.Logger.Error(err)
			return
		}
		err = tmpl.Execute(w, all)
		if err != nil {
			s.Logger.Error(err)
			return
		}
		s.Logger.Info("handle /getAllUsers GET")

	}
}

func (ctrl *Controller) UpdateUser(s *APIserver) http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		NewLogin string `json:"newLogin"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.Logger.Error(err)
			return
		}

		err := s.storage.User().ChangeLoginByEmail(req.NewLogin, req.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.Logger.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Возвращаем новый логин, чтобы обновить таблицу на клиенте
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"email":    req.Email,
			"newLogin": req.NewLogin,
		})
		s.Logger.Info("User login updated successfully")
	}
}

func (ctrl *Controller) DeleteUser(s *APIserver) http.HandlerFunc {
	type Request struct {
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.Logger.Error(err)
			return
		}

		err := s.storage.User().DeleteByEmail(req.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.Logger.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"msg":    "user deleted",
		})
		s.Logger.Info("User deleted successfully")
	}
}
