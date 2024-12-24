package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"

	storage "github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type ControllerUser struct {
	storage *storage.Storage
	session sessions.Store
	logger  *logrus.Logger
}

func NewControllerUser(storage *storage.Storage, session sessions.Store, logger *logrus.Logger) *ControllerUser {
	return &ControllerUser{
		storage: storage,
		session: session,
		logger:  logger,
	}
}

func (ctrl *ControllerUser) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := ctrl.storage.User().GetAll(20)

		if err != nil {
			Error(w, r, http.StatusNotFound, err)
			ctrl.logger.Error(err)
			return
		}
		w.Header().Set("Content-Type", "http")
		tmpl, err := template.ParseFiles("./templates/users.html")
		if err != nil {
			ctrl.logger.Error(err)
			return
		}
		err = tmpl.Execute(w, all)
		if err != nil {
			ctrl.logger.Error(err)
			return
		}
		ctrl.logger.Info("handle /getAllUsers GET")

	}
}

func (ctrl *ControllerUser) UpdateUser() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		NewLogin string `json:"newLogin"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctrl.logger.Error(err)
			return
		}

		err := ctrl.storage.User().ChangeLoginByEmail(req.NewLogin, req.Email)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Возвращаем новый логин, чтобы обновить таблицу на клиенте
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"email":    req.Email,
			"newLogin": req.NewLogin,
		})
		ctrl.logger.Info("User login updated successfully")
	}
}

func (ctrl *ControllerUser) DeleteUser() http.HandlerFunc {
	type Request struct {
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Error(w, r, http.StatusBadRequest, err)
			ctrl.logger.Error(err)
			return
		}

		err := ctrl.storage.User().DeleteByEmail(req.Email)
		if err != nil {
			Error(w, r, http.StatusBadRequest, err)

			ctrl.logger.Error(err)

			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"msg":    "user deleted",
		})
		ctrl.logger.Info("User deleted successfully")
	}
}

func (ctrl *ControllerUser) FindUser() http.HandlerFunc {
	type Request struct {
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		u, err := ctrl.storage.User().FindByEmail(req.Email)
		if err != nil {
			Error(w, r, http.StatusBadRequest, err)
			ctrl.logger.Warn("unhandle /findUser POST", err)
			return
		}
		w.Header().Set("Content-Type", "appication/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":    strconv.Itoa(u.ID),
			"email": u.Email,
			"login": u.Login,
		})
		ctrl.logger.Info("handle /findUser POST")

	}
}
