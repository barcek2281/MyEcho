package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"text/template"

	"github.com/barcek2281/MyEcho/internal/app/model"
)

const SessionName = "MyEcho"

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
		s.Respond(w, r, 201, map[string]interface{}{
			"posts": []map[string]string{
				{
					"content":    "wadwa",
					"author":     "wadw",
					"created_at": "wadwa",
				},
			},
		})

		s.Logger.Info("handle /getPage", r.URL)
	}

}
