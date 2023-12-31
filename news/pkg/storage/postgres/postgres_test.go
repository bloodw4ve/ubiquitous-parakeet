package postgres

import (
	"APIGateway/news/pkg/storage"
	"testing"
)

const dbURL = "postgres://postgres:PASSWORD@localhost:5432/news" // write your postgresdb password

func Test_postgres(t *testing.T) {
	db, err := New(dbURL)
	if err != nil {
		t.Fatal(err)
	}
	testCase := []storage.Post{{ID: 1,
		Title:   "Test Title",
		Content: "Test Content",
		PubTime: 0,
		Link:    "Test Link"}}
	err = db.PostMany(testCase)
	if err != nil {
		t.Fatal(err)
	}
	news, err := db.GetPosts(2, 0)
	if err != nil {
		t.Fatal(err)
	}
	const wantLen = 1
	if len(news) < wantLen {
		t.Fatalf("Got %d records, Wanted %d", len(news), wantLen)
	}
	_, _, err = db.PostSearch("", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.PostDetail(2)
	if err != nil {
		t.Fatal(err)
	}
	err = db.UpdatePost(testCase[0])
	if err != nil {
		t.Fatal(err)
	}
	err = db.DeletePost(testCase[0])
	if err != nil {
		t.Fatal(err)
	}
}
