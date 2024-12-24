package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"text/template"

	"golang.org/x/time/rate"

	"github.com/barcek2281/MyEcho/internal/app/model"
	storage "github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type ctxKey int8

const (
	SessionName       = "MyEcho"
	pageNumberDefault = 5
	ctxKeyUser        = 1
)

var (
	controllerPost ControllerPost

	errCantBeHere     = errors.New("you not suppose to be here")
	errSessionTimeOut = errors.New("your session time out")
	errTooManyRequest = errors.New("Too many request dude")

	limiter = rate.NewLimiter(1, 3)
)

type ControllerPost struct {
	storage *storage.Storage
	session sessions.Store
	logger  *logrus.Logger
}

func NewControllerPost(storage *storage.Storage, session sessions.Store, logger *logrus.Logger) *ControllerPost {
	return &ControllerPost{
		storage: storage,
		session: session,
		logger:  logger,
	}
}

func (ctrl *ControllerPost) CreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)

		tmpl, err := template.ParseFiles("./templates/post.html")
		if err != nil {
			ctrl.logger.Error(err)
			return
		}

		err = tmpl.Execute(w, u)
		if err != nil {
			ctrl.logger.Warn("cannot execute template", err)
			Error(w, r, 404, errSessionTimeOut)
			return
		}

		ctrl.logger.Info("Handle /createPost GET")
	}
}

// POST createPost/
func (ctrl *ControllerPost) CreatePostReal() http.HandlerFunc {
	type Request struct {
		Content string `json:"content"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// limit (1, 3)
		if !limiter.Allow() {
			Error(w, r, http.StatusTooManyRequests, errTooManyRequest)
			return
		}

		session, err := ctrl.session.Get(r, sessionName)
		if err != nil {
			Error(w, r, 404, errCantBeHere)
			ctrl.logger.Warn("anon cant be here")
			return
		}

		user_id, ok := session.Values["user_id"].(int)
		if !ok {
			Error(w, r, 404, errSessionTimeOut)
			ctrl.logger.Warn(err)
			return
		}
		req := Request{}
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			Error(w, r, 404, err)
			ctrl.logger.Warn(err)
			return
		}

		post := &model.Post{
			User_id: user_id,
			Content: req.Content,
		}
		if err := ctrl.storage.Post().Create(post); err != nil {
			Error(w, r, http.StatusBadRequest, err)
			ctrl.logger.Warn(err)
			return
		}

		Respond(w, r, http.StatusCreated, map[string]string{"status": "Succesfully, created post"})
		ctrl.logger.Info("handle /createPost POST")
	}
}

func (ctrl *ControllerPost) GetPost() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Создаем ответ с кодом 201 и передаем JSON
		login := r.URL.Query().Get("author")
		sortDate := r.URL.Query().Get("sort")
		if sortDate != "ASC" && sortDate != "DESC" {
			sortDate = "DESC"
		}
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page <= 0 {
			page = 1
		}

		posts, err := ctrl.storage.Post().GetAllWithAuthors(login, sortDate, pageNumberDefault, (page-1)*pageNumberDefault)
		if err != nil {
			ctrl.logger.Warn(err)
			Error(w, r, 504, err)
			return
		}
		res_posts := make([]map[string]string, 0)
		for _, post := range posts {
			res_posts = append(res_posts, map[string]string{
				"content":    post.Content,
				"author":     post.Author,
				"created_at": post.ConverDateToString(),
			})
		}
		Respond(w, r, http.StatusAccepted, map[string]interface{}{
			"posts": res_posts,
		})
		ctrl.logger.Info("handle /getPost ", r.URL)
	}

}
