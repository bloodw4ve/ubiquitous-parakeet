package rss

import (
	"APIGateway/news/pkg/storage"
	"testing"
)

func TestGetRss(t *testing.T) {
	news, err := GetRss("https://habr.com/ru/rss/hub/go/all/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(news) == 0 {
		t.Fatal("Couldn't get RSS-feed")
	}
}

func TestGoNews(t *testing.T) {
	chPost := make(chan []storage.Post)
	chErr := make(chan error)
	type args struct {
		configURL string
		chPost    chan []storage.Post
		chErr     chan error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{configURL: "./test_config.json", chPost: chPost, chErr: chErr},
			wantErr: false,
		},
		{
			name:    "null config",
			args:    args{configURL: "", chPost: chPost, chErr: chErr},
			wantErr: true,
		},
		{
			name:    "invalid config",
			args:    args{configURL: "./test_config_invalid.json", chPost: chPost, chErr: chErr},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GoNews(tt.args.configURL, tt.args.chPost, tt.args.chErr); (err != nil) != tt.wantErr {
				t.Errorf("GoNews() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
