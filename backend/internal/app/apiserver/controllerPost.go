package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"text/template"

	"github.com/barcek2281/MyEcho/internal/app/model"
)

const (
	SessionName       = "MyEcho"
	pageNumberDefault = 5
)

var (
	controllerPost ControllerPost

	errCantBeHere     = errors.New("you not suppose to be here")
	errSessionTimeOut = errors.New("your session time out")
)

type ControllerPost struct{}

func (ctrl *ControllerPost) CreatePost(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.Session.Get(r, sessionName)
		if err != nil {
			s.Error(w, r, 404, errCantBeHere)
			s.Logger.Warn("anon cant be here")
			return
		}

		id, ok := session.Values["user_id"].(int)
		if !ok {
			s.Error(w, r, 404, errSessionTimeOut)
			s.Logger.Warn(err)
			return
		}

		u := &model.User{}
		u, err = s.storage.User().FindById(id)
		if err != nil {
			s.Logger.Error("WTF how it happen?", err)
			s.Error(w, r, 404, errSessionTimeOut)
			return
		}

		tmpl, err := template.ParseFiles("./templates/post.html")
		if err != nil {
			s.Logger.Error(err)
			return
		}

		err = tmpl.Execute(w, u)
		if err != nil {
			s.Logger.Warn("cannot execute template", err)
			s.Error(w, r, 404, errSessionTimeOut)
		}

		s.Logger.Info("Handle /createPost GET")
	}
}

func (ctrl *ControllerPost) CreatePostReal(s *server) http.HandlerFunc {
	type Request struct {
		Content string `json:"content"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.Session.Get(r, sessionName)
		if err != nil {
			s.Error(w, r, 404, errCantBeHere)
			s.Logger.Warn("anon cant be here")
			return
		}

		user_id, ok := session.Values["user_id"].(int)
		if !ok {
			s.Error(w, r, 404, errSessionTimeOut)
			s.Logger.Warn(err)
			return
		}
		req := Request{}
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.Error(w, r, 404, err)
			s.Logger.Warn(err)
			return
		}

		post := &model.Post{
			User_id: user_id,
			Content: req.Content,
		}
		if err := s.storage.Post().Create(post); err != nil {
			s.Error(w, r, http.StatusBadRequest, err)
			s.Logger.Warn(err)
			return
		}

		s.Respond(w, r, http.StatusCreated, map[string]string{"status": "Succesfully, created post"})
		s.Logger.Info("handle /createPost POST")
	}
}

func (ctrl *ControllerPost) GetPost(s *server) http.HandlerFunc {

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

		posts, err := s.storage.Post().GetAllWithAuthors(login, sortDate, pageNumberDefault, (page-1)*pageNumberDefault)
		if err != nil {
			s.Logger.Warn(err)
			s.Error(w, r, 504, err)
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
		s.Respond(w, r, http.StatusAccepted, map[string]interface{}{"posts": res_posts})
		s.Logger.Info("handle /getPost ", r.URL)
	}

}
