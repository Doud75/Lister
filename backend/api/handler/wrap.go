package handler

import (
	"net/http"
	"setlist/api/apierror"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func Wrap(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if appErr := asAppError(err); appErr != nil {
				writeAppError(w, appErr)
			} else {
				writeAppError(w, apierror.InternalError(err.Error()))
			}
		}
	}
}
