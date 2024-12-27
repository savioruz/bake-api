package error

import (
	"encoding/json"
	"net/http"

	"github.com/savioruz/bake/internal/domain/model"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := model.NewErrorResponse(err)
	json.NewEncoder(w).Encode(response)
}
