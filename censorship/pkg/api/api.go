package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type API struct {
	Rout *mux.Router
}

func New() *API {
	api := API{
		Rout: mux.NewRouter(),
	}
	api.endpoints()
	return &api
}

func (api *API) endpoints() {
	api.Rout.HandleFunc("/comments/add", api.addCommentHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (api *API) Router() *mux.Router {
	return api.Rout
}

func (api *API) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	text := struct {
		Content string
	}{}
	err := json.NewDecoder(r.Body).Decode(&text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	banlist := []string{
		"qwerty",
		"asdfgh",
		"zxcvbn",
	}
	for _, banWord := range banlist {
		matched, err := regexp.MatchString(banWord, text.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if matched {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
