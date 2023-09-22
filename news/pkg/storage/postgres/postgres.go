package postgres

import (
	"APIGateway/news/pkg/storage"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

// connect to db
func New(conn string) (*Store, error) {
	log.Println("Connecting to news database...")
	for {
		_, err := pgxpool.Connect(context.Background(), conn)
		if err == nil {
			break
		}
	}
	db, err := pgxpool.Connect(context.Background(), conn)
	if err != nil {
		return nil, err
	}
	log.Println("Connected")
	return &Store{db: db}, nil
}

// PostSearch
func (s *Store) PostSearch(pattern string, limit, offset int) ([]storage.Post, storage.Pagination, error) {
	pattern = "%" + pattern + "%"

	pagination := storage.Pagination{
		Page:  offset/limit + 1,
		Limit: limit,
	}
	row := s.db.QueryRow(context.Background(), "SELECT count(*) FROM news WHERE title ILIKE $1;", pattern)
	err := row.Scan(&pagination.PageNum)

	if pagination.PageNum%limit > 0 {
		pagination.PageNum = pagination.PageNum/limit + 1
	} else {
		pagination.PageNum /= limit
	}
	if err != nil {
		return nil, storage.Pagination{}, err
	}
	rows, err := s.db.Query(context.Background(), "SELECT * FROM news WHERE title ILIKE $1 ORDER BY pubtime DESC LIMIT $2 OFFSET $3;", pattern, limit, offset)
	if err != nil {
		return nil, storage.Pagination{}, err
	}
	defer rows.Close()
	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.PubTime, &p.Link)
		if err != nil {
			return nil, storage.Pagination{}, err
		}
		posts = append(posts, p)
	}
	return posts, pagination, rows.Err()
}

// PostDetail
func (s *Store) PostDetail(id int) (storage.Post, error) {
	row := s.db.QueryRow(context.Background(), "SELECT * FROM news WHERE id = $1", id)
	var post storage.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
	if err != nil {
		return storage.Post{}, err
	}
	return post, nil
}

// GetPosts
func (s *Store) GetPosts(limit, offset int) ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM news ORDER BY pubtime DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.PubTime, &p.Link)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// AddPost
func (s *Store) AddPost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), "INSERT INTO news (title, content, pubtime, link) VALUES ($1, $2, $3, $4);", post.Title, post.Content, post.PubTime, post.Link)
	if err != nil {
		return err
	}
	return nil
}

// PostMany
func (s *Store) PostMany(posts []storage.Post) error {
	for _, post := range posts {
		err := s.AddPost(post)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdatePost
func (s *Store) UpdatePost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), "UPDATE news "+"SET title = $1, "+"content = $2, "+"pubtime = $3, "+"link = $4 "+"WHERE id = $5", post.Title, post.Content, post.PubTime, post.Link, post.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeletePost
func (s *Store) DeletePost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), "DELETE FROM news WHERE id = $1;", post.ID)
	if err != nil {
		return err
	}
	return nil
}

// Close
func (s *Store) Close() {
	s.db.Close()
}
