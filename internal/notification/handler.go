package notification

import (
	"encoding/json"
	"net/http"
)

type NotificationHandler struct {
	service *NotificationService
}

func NewNotificationHandler(service *NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) Notify(w http.ResponseWriter, r *http.Request) {
	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	notification, err := h.service.Notify(req)
	if err != nil {
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

func (h *NotificationHandler) FindNotificationByID(w http.ResponseWriter, r *http.Request) {
	notificationID := r.URL.Query().Get("notification_id")
	if notificationID == "" {
		http.Error(w, "Missing notification_id query parameter", http.StatusBadRequest)
		return
	}

	notification, err := h.service.FindNotificationByID(notificationID)
	if err != nil {
		http.Error(w, "Failed to fetch notification", http.StatusInternalServerError)
		return
	}
	if notification == nil {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}
