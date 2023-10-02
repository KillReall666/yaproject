package html

import (
	"fmt"
	"net/http"
)

type metricsHTML interface {
	PrintForHTML() string
}

type Handler struct {
	metricsHTML metricsHTML
}

func NewHTMLHandler(s metricsHTML) *Handler {
	return &Handler{
		metricsHTML: s,
	}
}

func (h *Handler) HTMLOutput(w http.ResponseWriter, r *http.Request) {
	htmlPage := h.metricsHTML.PrintForHTML()
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Encoding", "gzip")

	fmt.Fprint(w, htmlPage)

}
