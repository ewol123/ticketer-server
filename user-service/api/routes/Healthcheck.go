package routes

import "net/http"

// Get : check service health
func (h *handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	setupResponse(w, contentType, []byte{}, http.StatusOK)
}
