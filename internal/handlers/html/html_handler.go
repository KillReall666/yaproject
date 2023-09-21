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

func NewHtmlHandler(s metricsHTML) *Handler {
	return &Handler{
		metricsHTML: s,
	}
}

func (h *Handler) HTMLOutput(w http.ResponseWriter, r *http.Request) {
	htmlPage := h.metricsHTML.PrintForHTML()
	fmt.Fprint(w, htmlPage)

}
