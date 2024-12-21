package storage

import (
	

	"github.com/barcek2281/MyEcho/internal/app/model"
)

type PostRepository struct {
	storage *Storage
}

func (p *PostRepository) Create(post *model.Post) error {
	if err := post.Validate(); err != nil {
		return err
	}
	err := p.storage.db.QueryRow("INSERT INTO posts (user_id, content) VALUES ($1, $2) RETURNING id", post.User_id, post.Content).Scan(&post.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostRepository) GetAll(limit int) ([]*model.Post, error) {
	rows, err := p.storage.db.Query("SELECT user_id, content FROM posts ORDER BY created_at LIMIT $1", limit)
	if err != nil {
		return nil, err
	}

	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(&post.User_id, &post.Content); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *PostRepository) GetAllWithAuthors(login, sortDate string, limit, offset int) ([]*model.Post, error) {
	// запрос
	query := "SELECT posts.content, posts.user_id, users.login, posts.created_at FROM posts JOIN users ON posts.user_id = users.id "

	// имя пользователя
	if login != "" {
		query += "WHERE users.login = " + "'" + login + "'" 
	}

	query += " ORDER BY posts.created_at "
	if sortDate == "DESC" {
		query += sortDate
	} else if sortDate == "ASC" {
		query += sortDate
	} else {
		query += "DESC"
	}
	rows, err := p.storage.db.Query(query + " LIMIT $1", limit)
	
	if err != nil {
		return nil, err
	}

	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(&post.Content, &post.User_id, &post.Author, &post.Created_at); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
