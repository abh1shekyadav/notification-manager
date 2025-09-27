package notification

import (
	"encoding/json"
	"net/http"
)

type NotificationHandler struct {
	service *NotificationService
}

func NewNotificationHandler(notificationService *NotificationService) *NotificationHandler {
	return &NotificationHandler{service: notificationService}
}

func (h *NotificationHandler) Notify(w http.ResponseWriter, r *http.Request) {
	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	notification, err := h.service.Notify(req)
	if err != nil {
		http.Error(w, "failed to save notification", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}
