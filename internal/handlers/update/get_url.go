package update

import (
	"net/http"
	"strings"
)

func getURL(r *http.Request) []string {
	url := r.URL.String()

	urlWithoutPref, err := strings.CutPrefix(url, "/")
	if !err {
		panic(err)
	}

	requestString := strings.Split(urlWithoutPref, "/")

	return requestString
}
