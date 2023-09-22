package memdb

import "APIGateway/news/pkg/storage"

// data storage
type Store struct{}

// storage object construct
func New() *Store {
	return new(Store)
}

func (s *Store) PostSearch(keyWord string, limit, offset int) ([]storage.Post, storage.Pagination, error) {
	return nil, storage.Pagination{}, nil
}

func (s *Store) PostDetail(id int) (storage.Post, error) {
	return storage.Post{}, nil
}

func (s *Store) Posts(limit, offset int) ([]storage.Post, error) {
	return posts[0:0], nil
}

func (s *Store) AddPost(storage.Post) error {
	return nil
}

func (s *Store) PostMany([]storage.Post) error {
	return nil
}

func (s *Store) UpdatePost(storage.Post) error {
	return nil
}

func (s *Store) DeletePost(storage.Post) error {
	return nil
}

func (s *Store) Close() {
}

var posts = []storage.Post{
	{
		ID:      1,
		Title:   "Test_Title_1",
		Content: "Test_Content_1",
		PubTime: 0,
		Link:    "Test_Link_1"},
	{
		ID:      2,
		Title:   "Test_Title_2",
		Content: "Test_Content_2",
		PubTime: 0,
		Link:    "Test_Link_2"},
	{
		ID:      3,
		Title:   "Test_Title_3",
		Content: "Test_Content_3",
		PubTime: 0,
		Link:    "Test_Link_3"},
	{
		ID:      4,
		Title:   "Test_Title_4",
		Content: "Test_Content_4",
		PubTime: 0,
		Link:    "Test_Link_4"},
}
