package api

import (
	"encoding/json"
	"net/http"
)

func (h Handlers) savePushSubscription(w http.ResponseWriter, req *http.Request) {
	var payload PushSubscriptionRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		h.errorHandler.Write(w, ErrValidation)
		return
	}
	if payload.Endpoint == "" {
		h.errorHandler.Write(w, ErrValidation)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "subscription stored"})
}
