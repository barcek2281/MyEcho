package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/barcek2281/MyEcho/internal/app/mail"
	"github.com/barcek2281/MyEcho/internal/app/model"
	storage "github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/barcek2281/MyEcho/pkg/utils"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName      = "MyEcho"
	sessionAdmin     = "IsAdmin"
	roleUser         = "user"
	roleUnauthorized = "unauthorized"
	roleAdmin        = "admin"
)

var (
	errIncorrectPasswordOrEmail = errors.New("Incorrect password or email")
	errYouDontHaveUsers         = errors.New("YOU DONT HAVE USERS DUMBASS")
)

type Controller struct {
	storage *storage.Storage
	session sessions.Store
	logger  *logrus.Logger
	sender  *mail.Sender
}

func NewController(storage *storage.Storage, session sessions.Store, logger *logrus.Logger, sender *mail.Sender) *Controller {
	return &Controller{
		storage: storage,
		session: session,
		logger:  logger,
		sender:  sender,
	}
}

func (ctrl *Controller) MainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			ctrl.logger.Error(err)
			return
		}
		var user *model.User = nil

		session, err := ctrl.session.Get(r, sessionName)
		if err != nil {
			ctrl.logger.Info("no session")
		}
		if session != nil {
			userID, ok := session.Values["user_id"].(int)
			if !ok {
				ctrl.logger.Warn("session timeout!", err)
			} else {
				user, err = ctrl.storage.User().FindById(userID)
				if err != nil {
					ctrl.logger.Warn("warn lol )", err)
				}
			}
		}

		data := map[string]interface{}{
			"user": user,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			ctrl.logger.Error(err)
			return
		}

		ctrl.logger.Info("handle MainPage GET")
	}
}

func (ctrl *Controller) HandleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"status":  "OK",
			"message": "Hello World!",
		}
		ctrl.logger.Info("handle /hello")

		// Print query
		prettyJson, err := json.MarshalIndent(r.URL.Query(), "", "  ")
		if err != nil {
			ctrl.logger.Error(err)
		}
		fmt.Println(string(prettyJson))

		// receive accept
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			ctrl.logger.Error(err)
			http.Error(w, "cannot write json file", http.StatusInternalServerError)
		}
	}
}

func (ctrl *Controller) HandleHelloPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := map[string]string{}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "not ok", "msg": "we got a cringe/error message"})
			return
		}

		if _, ok := req["msg"]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "not ok", "msg": "we didnt get a message "})
			return
		}

		fmt.Println(req["msg"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "msg": "we got the message"})

		ctrl.logger.Info("handle /hello POST")
	}
}

func (ctrl *Controller) RegisterPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/register.html")
		if err != nil {
			ctrl.logger.Error(err)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			ctrl.logger.Error(err)
			return
		}
		ctrl.logger.Info("handle /register GET")
	}
}

func (ctrl *Controller) RegisterUser() http.HandlerFunc {
	type Request struct {
		Login    string `json:"login"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			ctrl.logger.Error(err)
			return
		}

		u := model.User{
			Email:    req.Email,
			Password: req.Password,
			Login:    req.Login,
		}

		if err := ctrl.storage.User().Create(&u); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			ctrl.logger.Warn(err)
			return
		}

		session, err := ctrl.session.Get(r, sessionName)
		if err != nil {
			utils.Error(w, r, 404, err)
			ctrl.logger.Warn(err)
			return
		}

		session.Values["user_id"] = u.ID
		// session.Values["role"] = roleUnauthorized

		if err := ctrl.session.Save(r, w, session); err != nil {
			utils.Error(w, r, 404, err)
			return
		}

		randomInt := rand.Intn(10000)
		bar := model.Barcode{
			User_id: u.ID,
			Barcode: randomInt,
		}
		if err := ctrl.storage.Barcode().Create(&bar); err != nil {
			ctrl.logger.Warn(err)
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		go ctrl.sender.SendToPerson("Your code", "your code: "+strconv.Itoa(randomInt), []string{u.Email})
		utils.Response(w, r, http.StatusCreated, map[string]string{"status": "Succesfully, created user"})
		ctrl.logger.Info("handle /register POST")
	}
}

func (ctrl *Controller) EmailVerifyPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmp, err := template.ParseFiles("./templates/email_verification.html")
		if err != nil {
			ctrl.logger.Warn(err)
			return
		}
		tmp.Execute(w, nil)
		ctrl.logger.Info("handle /register/verify GET")
	}
}

func (ctrl *Controller) EmailVerifyUser() http.HandlerFunc {
	type Request struct {
		Barcode int `json:"barcode"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			ctrl.logger.Warn(err)
			utils.Error(w, r, http.StatusBadGateway, err)
			return
		}
		session, err := ctrl.session.Get(r, SessionName)
		if err != nil {
			ctrl.logger.Warn(err)
			utils.Error(w, r, http.StatusBadGateway, err)
			return
		}

		user_id := session.Values["user_id"].(int)
		barcode, err := ctrl.storage.Barcode().FindByUserId(user_id)
		if barcode.Barcode != req.Barcode {
			ctrl.logger.Warn(err)
			utils.Error(w, r, http.StatusBadGateway, errIncorrectPasswordOrEmail)
			return
		}
		session.Values["role"] = roleUser
		if err := ctrl.session.Save(r, w, session); err != nil {
			ctrl.logger.Warn(err)
			utils.Error(w, r, 404, err)
			return
		}
		utils.Response(w, r, 200, nil)
		ctrl.logger.Info("handle /register/verify POST")
	}
}

func (ctrl *Controller) LoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/login.html")
		if err != nil {
			ctrl.logger.Warn(err)
			// ctrl.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			ctrl.logger.Error(err)
			return
		}
		ctrl.logger.Info("Handle /login GET")
	}
}

func (ctrl *Controller) LoginUser() http.HandlerFunc {
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

		u, err := ctrl.storage.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			utils.Error(w, r, 404, errIncorrectPasswordOrEmail)
			return
		}

		session, err := ctrl.session.New(r, sessionName)
		if err != nil {
			utils.Error(w, r, 404, err)
			return
		}
		session.Values["user_id"] = u.ID
		session.Values["role"] = roleUser
		err = ctrl.session.Save(r, w, session)

		w.Header().Set("Access-Control-Allow-Origin", "http://example.com") // Замени на домен клиента
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		ctrl.logger.Info("handle /login POST")
		// http.Redirect(w, r, "/", http.StatusSeeOther)
		utils.Response(w, r, 201, nil)
	}
}

func (ctrl *Controller) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение сессии
		session, err := ctrl.session.Get(r, sessionName)
		if err != nil {
			ctrl.logger.Warn("Failed to get session: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Удаление данных из сессии
		session.Options.MaxAge = -1
		err = session.Save(r, w)
		if err != nil {
			ctrl.logger.Warn("Failed to delete session: ", err)
			// ctrl.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		ctrl.logger.Info("handle /logout ANY")
	}
}

func (ctrl *Controller) SupportPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/support.html")
		if err != nil {
			ctrl.logger.Error(err)
			return
		}

		tmpl.Execute(w, nil)
		ctrl.logger.Info("handle support/ GET")
	}
}

func (ctrl *Controller) SupportUser() http.HandlerFunc {
	type Request struct {
		TypeProblem string `json:"type"`
		Text        string `json:"text"`
		Filename    string `json:"filename"`
		File        string `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiterSupport.Allow() {
			utils.Error(w, r, http.StatusTooManyRequests, errTooManyRequest)
			return
		}
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			ctrl.logger.Error(err)
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}
		email := "email: anonymous"

		session, err := ctrl.session.Get(r, sessionName)
		if err == nil {
			id, ok := session.Values["user_id"].(int)
			if ok {
				u, err := ctrl.storage.User().FindById(id)
				if err == nil {
					email = "email: " + u.Email
				}
			}
		}

		err = ctrl.sender.SendToSupport(req.TypeProblem, req.Text, email, req.Filename, &req.File)
		if err != nil {
			ctrl.logger.Error(err)
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}
		ctrl.logger.Info("handle support/ POST")
	}
}
