package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type API struct {
	Rt *mux.Router
}

type ResponseDetailed struct {
	NewsDetailed struct {
		ID      int    `json:"ID"`
		Title   string `json:"Title"`
		Content string `json:"Content"`
		PubTime int    `json:"PubTime"`
		Link    string `json:"Link"`
	} `json:"NewsDetailed"`
	Comments []struct {
		ID              int    `json:"ID"`
		NewsID          int    `json:"newsID"`
		ParentCommentID int    `json:"parentCommentID"`
		Content         string `json:"content"`
		PubTime         int    `json:"pubTime"`
	} `json:"Comments"`
}

const limit = 10 // one-page post limit
const newsService = "http://localhost:8080"
const commentService = "http://localhost:8081"
const censorAddr = "http://localhost:8082"

// creates API
func New() *API {
	api := API{
		Rt: mux.NewRouter(),
	}
	api.endpoints()
	return &api
}

// API handlers
func (api *API) endpoints() {
	api.Rt.HandleFunc("/news", api.newsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.Rt.HandleFunc("/news/latest", api.newsLatestHandler).Methods(http.MethodGet, http.MethodOptions)
	api.Rt.HandleFunc("/comments/add", api.commentHandler).Methods(http.MethodPost, http.MethodOptions)
	api.Rt.HandleFunc("/news/search", api.newsDetailedHandler).Methods(http.MethodGet, http.MethodOptions)
}

func MakeRequest(req *http.Request, method, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	r.Header = req.Header
	return client.Do(r)
}

// request router
func (api *API) Router() *mux.Router {
	return api.Rt
}

func (api *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	sParam := r.URL.Query().Get("s")

	resp, err := MakeRequest(r, http.MethodGet, newsService+"/news"+"?page="+pageParam+"&"+"s="+sParam, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (api *API) newsLatestHandler(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	resp, err := MakeRequest(r, http.MethodGet, newsService+"/news/latest"+"?page="+pageParam, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
func (api *API) commentHandler(w http.ResponseWriter, r *http.Request) {
	bodybytes, _ := io.ReadAll(r.Body)
	r.Body.Close()
	Body1 := io.NopCloser(bytes.NewBuffer(bodybytes))
	Body2 := io.NopCloser(bytes.NewBuffer(bodybytes))
	respCensor, err := MakeRequest(r, http.MethodPost, censorAddr+"/comments/add", Body1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if respCensor.StatusCode != 200 {
		http.Error(w, "incorrect comment contents", respCensor.StatusCode)
		return
	}
	resp, err := MakeRequest(r, http.MethodPost, commentService+"/comments/add", Body2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (api *API) newsDetailedHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "no search parameters found", http.StatusBadRequest)
		return
	}
	chNews := make(chan *http.Response, 2)
	chComments := make(chan *http.Response, 2)
	chErr := make(chan error, 2)
	var response ResponseDetailed
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		respComments, err := MakeRequest(r, http.MethodGet, commentService+"/comments"+"?news_id="+idParam, nil)
		chErr <- err
		chComments <- respComments
	}()
	wg.Wait()
	close(chErr)

	for err := range chErr {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
block:
	for {
		select {
		case respNews := <-chNews:
			body, err := io.ReadAll(respNews.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.Unmarshal(body, &response.NewsDetailed)
		case respComment := <-chComments:
			body, err := io.ReadAll(respComment.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.Unmarshal(body, &response.Comments)
		default:
			break block
		}
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
