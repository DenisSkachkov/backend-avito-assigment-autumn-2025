package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/dto"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/pullrequest"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/user"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService user.UserService
	prService pullrequest.PullRequestService
}

func NewUserHandler(userService user.UserService, prService pullrequest.PullRequestService) *UserHandler {
    return &UserHandler{userService: userService, prService: prService}
}

func (h *UserHandler) RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/users/setIsActive", h.SetIsActive).Methods(http.MethodPost)
    r.HandleFunc("/users/getReview", h.GetReviewPRs).Methods(http.MethodGet)
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
    var body dto.UserDTO

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user, err := h.userService.SetActive(r.Context(), body.UserID, body.IsActive)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "user": user,
    })
}

func (h *UserHandler) GetReviewPRs(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    if userID == "" {
        http.Error(w, "user_id is required", http.StatusBadRequest)
        return
    }

    prs, err := h.prService.GetReview(r.Context(), userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "user_id":       userID,
        "pull_requests": prs,
    })
}
