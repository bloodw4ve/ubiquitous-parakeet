package rss

import (
	"APIGateway/news/pkg/storage"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
)

type Channel struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type config struct {
	Rss           []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

// get rss feed by url
func GetRss(url string) ([]storage.Post, error) {
	var c Channel
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	err = xml.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		return nil, err
	}
	// received data -> publication list
	var news []storage.Post
	for _, item := range c.Items {
		var p storage.Post
		p.Title = item.Title
		p.Content = item.Description
		p.Content = strip.StripTags(p.Content)
		p.Link = item.Link

		t, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			t, err = time.Parse(time.RFC1123Z, item.PubDate)
		}
		if err != nil {
			t, err = time.Parse("Mon, _2 Jan 2006 15:04:05 -0700", item.PubDate)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}
		news = append(news, p)
	}
	return news, nil
}

// Reads rss from config file using goroutines for each rss feed
func GoNews(configURL string, chPost chan<- []storage.Post, chErr chan<- error) error {
	file, err := os.Open(configURL)
	if err != nil {
		return err
	}
	var conf config
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		return err
	}
	log.Println("Watching rss feeds has begun...")
	for i, r := range conf.Rss {
		go func(r string, i int, chPost chan<- []storage.Post, chErr chan<- error) {
			for {
				log.Println("Goroutine is launched, ", i, " link: ", r)
				p, err := GetRss(r)
				if err != nil {
					chErr <- err
					continue
				}
				chPost <- p
				log.Println("Inserted posts from goroutine, ", i, " link: ", r)
				log.Println("Goroutine ", i, " is waiting for the next iteration")
				time.Sleep(time.Duration(conf.RequestPeriod) * time.Minute)
			}
		}(r, i, chPost, chErr)
	}
	return nil
}
