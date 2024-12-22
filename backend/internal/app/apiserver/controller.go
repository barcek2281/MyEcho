package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"

	"text/template"

	"github.com/barcek2281/MyEcho/internal/app/model"
	_ "github.com/sirupsen/logrus"
)

const (
	sessionName = "MyEcho"
)

var (
	controller                  Controller
	errIncorrectPasswordOrEmail = errors.New("incorrect password or email")
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) MainPage(s *server) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			s.Logger.Error(err)
			return
		}
		var user *model.User = nil

		session, err := s.Session.Get(r, sessionName)
		if err != nil {
			s.Logger.Info("no session")
		}

		if session != nil {
			userID, ok := session.Values["user_id"].(int)
			if !ok {
				s.Logger.Warn("session timeout!", err)
			} else {
				user, err = s.storage.User().FindById(userID)
				if err != nil {
					s.Logger.Warn("warn lol )", err)
				}
			}
		}

		data := map[string]interface{}{
			"user": user,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			s.Logger.Error(err)
			return
		}

		s.Logger.Info("handle MainPage GET")
	}
}

func (ctrl *Controller) handleHello(s *server) http.HandlerFunc {
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

func (ctrl *Controller) handleHelloPost(s *server) http.HandlerFunc {
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

		s.Logger.Info("handle /hello POST")
	}
}

func (ctrl *Controller) registerPage(s *server) http.HandlerFunc {
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

func (ctrl *Controller) registerUser(s *server) http.HandlerFunc {
	type Request struct {
		Login    string `json:"login"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.Error(w, r, http.StatusBadRequest, err)
			s.Logger.Error(err)
			return
		}

		u := model.User{
			Email:    req.Email,
			Password: req.Password,
			Login:    req.Login,
		}

		if err := s.storage.User().Create(&u); err != nil {
			s.Error(w, r, http.StatusBadRequest, err)
			s.Logger.Warn(err)
			return
		}

		session, err := s.Session.Get(r, sessionName)
		if err != nil {
			s.Error(w, r, 404, err)
			s.Logger.Warn(err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.Session.Save(r, w, session); err != nil {
			s.Error(w, r, 404, err)
			return
		}
		s.Respond(w, r, http.StatusCreated, map[string]string{"status": "Succesfully, created user"})

		s.Logger.Info("handle /register POST")
	}
}

func (ctrl *Controller) loginPage(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/login.html")
		if err != nil {
			s.Logger.Warn(err)
			s.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			s.Logger.Error(err)
			return
		}
		s.Logger.Info("Handle /login GET")
	}
}

func (ctrl *Controller) loginUser(s *server) http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.Logger.Warn(err)
			s.Error(w, r, 404, err)
			return
		}

		u, err := s.storage.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.Error(w, r, 404, errIncorrectPasswordOrEmail)
			return
		}
		session, err := s.Session.Get(r, sessionName)
		if err != nil {
			s.Error(w, r, 404, err)
			return
		}
		session.Values["user_id"] = u.ID
		s.Session.Save(r, w, session)

		s.Logger.Info("handle /login POST")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (ctrl *Controller) LogoutHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение сессии
		session, err := s.Session.Get(r, sessionName)
		if err != nil {
			s.Logger.Warn("Failed to get session: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Удаление данных из сессии
		session.Options.MaxAge = -1 // Устанавливаем MaxAge в -1 для удаления куки
		err = session.Save(r, w)
		if err != nil {
			s.Logger.Warn("Failed to delete session: ", err)
			s.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Перенаправление на главную страницу или страницу входа
		http.Redirect(w, r, "/", http.StatusSeeOther)
		s.Logger.Info("handle /logout ANY")
	}
}

func (ctrl *Controller) SupportPage(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/support.html")
		if err != nil {
			s.Logger.Error(err)
			return
		}

		tmpl.Execute(w, nil)
		s.Logger.Info("handle support/ GET")
	}
}

func (ctrl *Controller) SupportUser(s *server) http.HandlerFunc {
	type Request struct {
		TypeProblem string `json:"type"`
		Text        string `json:"text"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.Logger.Error(err)
			s.Error(w, r, http.StatusBadRequest, err)
			return
		}
		email := ""

		session, err := s.Session.Get(r, sessionName)
		if err == nil {
			id, ok := session.Values["user_id"].(int)
			if ok {
				u, err := s.storage.User().FindById(id)
				if err == nil {
					email = u.Email
				}
			}
		}

		auth := smtp.PlainAuth("hitler", s.Env.EmailTo, s.Env.EmailToPassword, "smtp.gmail.com")

		to := []string{s.Env.EmailTo}

		msg := "Subject: " + req.TypeProblem +
			"\r\n\n" +
			req.Text + "\r\n"

		if email != "" {
			msg += "email: " + email
		} else {
			msg += "email: anonymous"
		}

		err = smtp.SendMail("smtp.gmail.com:587", auth, "", to, []byte(msg))
		if err != nil {
			s.Logger.Warn(err)
			s.Error(w, r, http.StatusBadRequest, err)
			return
		}

		s.Respond(w, r, http.StatusAccepted, nil)
		s.Logger.Info("handle support/ POST")
	}
}
