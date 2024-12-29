package helper

import (
	"net/http"
	"strings"
)

func ParseParam(r *http.Request) string {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
