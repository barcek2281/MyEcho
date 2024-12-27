package apiserver

import (
	"encoding/json" // Для работы с JSON
	"errors"        // Для создания пользовательских ошибок
	"net/http"      // Для обработки HTTP-запросов
	"strconv"       // Для преобразования строк в числа

	"text/template" // Для работы с HTML-шаблонами

	"golang.org/x/time/rate" // Для ограничения частоты запросов

	"github.com/barcek2281/MyEcho/internal/app/model" // Импорт моделей приложения
)

const (
	SessionName       = "MyEcho" // Имя сессии
	pageNumberDefault = 5         // Количество записей на страницу
)

var (
	controllerPost ControllerPost // Контроллер для обработки запросов, связанных с постами

	errCantBeHere     = errors.New("you not suppose to be here") // Ошибка, если пользователь находится не там, где должен
	errSessionTimeOut = errors.New("your session time out")      // Ошибка, если сессия истекла
	errTooManyRequest = errors.New("Too many request dude")       // Ошибка, если превышен лимит запросов

	limiter = rate.NewLimiter(1, 3) // Лимитер: 1 запрос в секунду с буфером до 3 запросов
)

type ControllerPost struct{} // Структура контроллера постов

// Метод для обработки GET-запроса на создание поста
func (ctrl *ControllerPost) CreatePost(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.Session.Get(r, sessionName) // Получаем сессию пользователя
		if err != nil {
			s.Error(w, r, 404, errCantBeHere) // Ошибка: пользователь не авторизован
			s.Logger.Warn("anon cant be here") // Логируем предупреждение
			return
		}

		id, ok := session.Values["user_id"].(int) // Проверяем, есть ли user_id в сессии
		if !ok {
			s.Error(w, r, 404, errSessionTimeOut) // Ошибка: сессия истекла
			s.Logger.Warn(err)                   // Логируем ошибку
			return
		}

		u := &model.User{} // Создаем объект пользователя
		u, err = s.storage.User().FindById(id) // Ищем пользователя в базе по ID
		if err != nil {
			s.Logger.Error("WTF how it happen?", err) // Логируем ошибку поиска
			s.Error(w, r, 404, errSessionTimeOut)      // Ошибка: пользователь не найден
			return
		}

		tmpl, err := template.ParseFiles("./templates/post.html") // Загружаем HTML-шаблон
		if err != nil {
			s.Logger.Error(err) // Логируем ошибку загрузки шаблона
			return
		}

		err = tmpl.Execute(w, u) // Заполняем шаблон данными пользователя и отправляем ответ
		if err != nil {
			s.Logger.Warn("cannot execute template", err) // Логируем предупреждение
			s.Error(w, r, 404, errSessionTimeOut)          // Ошибка выполнения шаблона
		}

		s.Logger.Info("Handle /createPost GET") // Логируем успешную обработку запроса
	}
}

// Метод для обработки POST-запроса на создание реального поста
func (ctrl *ControllerPost) CreatePostReal(s *server) http.HandlerFunc {
	type Request struct {
		Content string `json:"content"` // Структура для хранения контента поста
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, можно ли выполнять запрос с учетом лимита
		if !limiter.Allow() {
			s.Error(w, r, http.StatusTooManyRequests, errTooManyRequest) // Ошибка: слишком много запросов
			return
		}

		session, err := s.Session.Get(r, sessionName) // Получаем сессию пользователя
		if err != nil {
			s.Error(w, r, 404, errCantBeHere) // Ошибка: пользователь не авторизован
			s.Logger.Warn("anon cant be here") // Логируем предупреждение
			return
		}

		user_id, ok := session.Values["user_id"].(int) // Извлекаем user_id из сессии
		if !ok {
			s.Error(w, r, 404, errSessionTimeOut) // Ошибка: сессия истекла
			s.Logger.Warn(err)                   // Логируем ошибку
			return
		}
		req := Request{} // Создаем объект для хранения данных запроса
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.Error(w, r, 404, err) // Ошибка: некорректный JSON в запросе
			s.Logger.Warn(err)      // Логируем ошибку
			return
		}

		post := &model.Post{
			User_id: user_id,       // ID пользователя
			Content: req.Content,   // Контент поста
		}
		if err := s.storage.Post().Create(post); err != nil {
			s.Error(w, r, http.StatusBadRequest, err) // Ошибка: пост не создан
			s.Logger.Warn(err)                        // Логируем ошибку
			return
		}

		s.Respond(w, r, http.StatusCreated, map[string]string{"status": "Succesfully, created post"}) // Успешный ответ
		s.Logger.Info("handle /createPost POST") // Логируем успешное выполнение запроса
	}
}

// Метод для получения списка постов
func (ctrl *ControllerPost) GetPost(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем параметры из URL
		login := r.URL.Query().Get("author")        // Фильтрация по автору
		sortDate := r.URL.Query().Get("sort")       // Сортировка (ASC/DESC)
		if sortDate != "ASC" && sortDate != "DESC" { // Если сортировка не задана, устанавливаем DESC
			sortDate = "DESC"
		}
		page, _ := strconv.Atoi(r.URL.Query().Get("page")) // Номер страницы
		if page <= 0 {
			page = 1 // Устанавливаем первую страницу, если номер некорректный
		}

		// Получаем список постов с учетом фильтров и пагинации
		posts, err := s.storage.Post().GetAllWithAuthors(login, sortDate, pageNumberDefault, (page-1)*pageNumberDefault)
		if err != nil {
			s.Logger.Warn(err) // Логируем ошибку получения постов
			s.Error(w, r, 504, err) // Ошибка: посты не получены
			return
		}
		res_posts := make([]map[string]string, 0) // Создаем массив для преобразования постов в JSON
		for _, post := range posts {
			res_posts = append(res_posts, map[string]string{
				"content":    post.Content,               // Контент поста
				"author":     post.Author,                // Автор поста
				"created_at": post.ConverDateToString(), // Дата создания поста
			})
		}
		s.Respond(w, r, http.StatusAccepted, map[string]interface{}{
			"posts": res_posts, // Отправляем массив постов в ответ
		})
		s.Logger.Info("handle /getPost ", r.URL) // Логируем успешное выполнение запроса
	}
}