package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"setlist/api/apierror"
)

func writeAppError(w http.ResponseWriter, appErr *apierror.AppError) {
	if !appErr.IsUserError {
		log.Printf("[ERROR][%s] %s", appErr.Code, appErr.Message)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatus)
	json.NewEncoder(w).Encode(map[string]string{
		"error": appErr.Message,
		"code":  appErr.Code,
	})
}
