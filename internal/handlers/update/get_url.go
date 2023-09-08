package update

import "net/http"

func getURL(r *http.Request) string {
	url := r.URL.String()
	return url
}
