package easy_test

import "net/http"

type JsonData struct {
	Nya string `json:"nya"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
