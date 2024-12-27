package helper

import (
	"net/http"
)

func ParseParam(r *http.Request) string {
	id := r.URL.Query().Get("id")
	return id
}
