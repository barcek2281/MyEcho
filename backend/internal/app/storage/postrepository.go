package storage

import "github.com/barcek2281/MyEcho/internal/app/model"

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

func (p *PostRepository) GetAllWithAuthors(limit int) ([]*model.Post, error) {
	rows, err := p.storage.db.Query("SELECT posts.content, posts.user_id, users.login FROM posts JOIN users ON posts.user_id = users.id ORDER BY posts.created_at LIMIT $1", limit)
	if err != nil {
		return nil, err
	}

	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(&post.Content, &post.User_id, &post.Author); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
