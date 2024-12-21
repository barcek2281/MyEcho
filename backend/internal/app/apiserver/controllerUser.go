package apiserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
)


type ControllerUser struct{}

var controllerUser ControllerUser

func (ctrl *ControllerUser) getAllUsers(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := s.storage.User().GetAll(20)

		if err != nil {
			s.Error(w, r, http.StatusNotFound, err)
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

func (ctrl *ControllerUser) UpdateUser(s *server) http.HandlerFunc {
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
			s.Error(w, r, http.StatusInternalServerError, err)
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

func (ctrl *ControllerUser) DeleteUser(s *server) http.HandlerFunc {
	type Request struct {
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.Error(w, r, http.StatusBadRequest, err)
			s.Logger.Error(err)
			return
		}

		err := s.storage.User().DeleteByEmail(req.Email)
		if err != nil {
			s.Error(w, r, http.StatusBadRequest, err)

			s.Logger.Error(err)

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

func (ctrl *ControllerUser) FindUser(s *server) http.HandlerFunc {
	type Request struct {
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		u, err := s.storage.User().FindByEmail(req.Email)
		if err != nil {
			s.Error(w, r, http.StatusBadRequest, err)
			s.Logger.Warn("unhandle /findUser POST", err)
			return
		}
		w.Header().Set("Content-Type", "appication/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":    strconv.Itoa(u.ID),
			"email": u.Email,
			"login": u.Login,
		})
		s.Logger.Info("handle /findUser POST")

	}
}
