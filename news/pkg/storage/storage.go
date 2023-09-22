package storage

// publication model
type Post struct {
	ID      int    `json:"ID,omitempty"`      // post id
	Title   string `json:"title,omitempty"`   // post title
	Content string `json:"content,omitempty"` // post content
	PubTime int64  `json:"pubTime,omitempty"` // post publication time
	Link    string `json:"link,omitempty"`    // post source link
}

type Pagination struct {
	PageNum int `json:"page_num,omitempty"`
	Page    int `json:"page,omitempty"`
	Limit   int `json:"limit,omitempty"`
}

// db interface
type Interface interface {
	PostSearch(keyWord string, limit, offset int) ([]Post, Pagination, error) // search
	PostDetail(id int) (Post, error)                                          // detailed
	GetPosts(limit, offset int) ([]Post, error)                               // get n posts
	AddPost(Post) error                                                       // add new post
	PostMany([]Post) error                                                    // add n posts // check
	UpdatePost(Post) error                                                    // update post
	DeletePost(Post) error                                                    // delete post(id)
	Close()                                                                   // closes connection to db
}
