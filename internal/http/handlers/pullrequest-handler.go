package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/dto"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/pullrequest"
	"github.com/gorilla/mux"
)

type PullRequestHandler struct {
	prService pullrequest.PullRequestService
}

func NewPullRequestHandler(prService pullrequest.PullRequestService) *PullRequestHandler {
    return &PullRequestHandler{prService: prService}
}

func (h *PullRequestHandler) RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/pullRequest/create", h.CreatePR).Methods(http.MethodPost)
    r.HandleFunc("/pullRequest/merge", h.MergePR).Methods(http.MethodPost)
    r.HandleFunc("/pullRequest/reassign", h.ReassignReviewer).Methods(http.MethodPost)
}

func (h *PullRequestHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
    var body dto.CreatePRDTO

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    pr, err := h.prService.Create(r.Context(), body.PullRequestID, body.AuthorID, body.PullRequestName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "pr": pr,
    })
}

func (h *PullRequestHandler) MergePR(w http.ResponseWriter, r *http.Request) {
    var body dto.MergePRDTO

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    pr, err := h.prService.Merge(r.Context(), body.PullRequestID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "pr": pr,
    })
}

func (h *PullRequestHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
    var body dto.ReassignReviewerDTO

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    pr, replacedBy, err := h.prService.ReassignReviewer(r.Context(), body.PullRequestID, body.OldUserID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusConflict)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "pr":          pr,
        "replaced_by": replacedBy,
    })
}

