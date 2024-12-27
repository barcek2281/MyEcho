package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/barcek2281/MyEcho/internal/app/model"
	storage "github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/barcek2281/MyEcho/internal/app/mail"
	"github.com/barcek2281/MyEcho/pkg/utils"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type ControllerUser struct {
	storage *storage.Storage
	session sessions.Store
	logger  *logrus.Logger
	sender  *mail.Sender
}

func NewControllerUser(storage *storage.Storage, session sessions.Store, logger *logrus.Logger, sender *mail.Sender) *ControllerUser {
	return &ControllerUser{
		storage: storage,
		session: session,
		logger:  logger,
		sender:  sender,
	}
}

func (ctrl *ControllerUser) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := ctrl.storage.User().GetAll(20)

		if err != nil {
			utils.Error(w, r, http.StatusNotFound, err)
			ctrl.logger.Error(err)
			return
		}
		w.Header().Set("Content-Type", "http")
		tmpl, err := template.ParseFiles("./templates/admin_panel.html")
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
			utils.Error(w, r, http.StatusInternalServerError, err)
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
			utils.Error(w, r, http.StatusBadRequest, err)
			ctrl.logger.Error(err)
			return
		}

		err := ctrl.storage.User().DeleteByEmail(req.Email)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
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
			utils.Error(w, r, http.StatusBadRequest, err)
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

func (ctrl *ControllerUser) SendMessageAdmin() http.HandlerFunc {
	type Request struct {
		Msg string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}
		users, err := ctrl.storage.User().GetAllWithoutLimit()
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		people := make([]string, 0)
		for _, user := range users {
			people = append(people, user.Email)
		}

		if len(people) <= 0 {
			utils.Error(w, r, 503, errYouDontHaveUsers)
			return
		}

		err = ctrl.sender.SendToEveryPerson("hello, dear users", req.Msg, people)
		if err != nil {
			fmt.Println(err)
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		utils.Response(w, r, 202, nil)
		ctrl.logger.Info("Handle /admin/sendMessage POST")
	}
}

func (ctrl *ControllerUser) AdminLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/admin_login.html")
		if err != nil {
			utils.Error(w, r, 504, err)
			return
		}

		tmpl.Execute(w, nil)
		ctrl.logger.Info("handle /admin/ GET")
	}
}

func (ctrl *ControllerUser) AdminLogin() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			ctrl.logger.Warn(err)
			utils.Error(w, r, 404, err)
			return
		}

		a, err := ctrl.storage.Admin().FindByEmail(req.Email)
		if err != nil || !a.ComparePassword(req.Password) {
			utils.Error(w, r, 404, errIncorrectPasswordOrEmail)
			return
		}

		session, err := ctrl.session.Get(r, sessionAdmin)
		session.Values["admin_id"] = a.ID
		ctrl.session.Save(r, w, session)
		http.Redirect(w, r, "/admin/users", 202)
		ctrl.logger.Info("handle /admin/login POST")
	}
}

func (ctrl *ControllerUser) AdminRegister() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		a := &model.Admin{
			Email:    req.Email,
			Name:     req.Name,
			Password: req.Password,
		}

		err := ctrl.storage.Admin().Create(a)
		if err != nil {
			utils.Error(w, r, http.StatusBadGateway, err)
		}

		ctrl.logger.Info("Handle /admin/register/ POST")
	}
}
