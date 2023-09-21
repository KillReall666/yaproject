package handlers

import (
	"net/http"
	"strings"
)

func GetURL(r *http.Request) []string {
	url := r.URL.String()

	urlWithoutPref, err := strings.CutPrefix(url, "/")
	if !err {
		panic(err)
	}

	requestString := strings.Split(urlWithoutPref, "/")

	return requestString
}
