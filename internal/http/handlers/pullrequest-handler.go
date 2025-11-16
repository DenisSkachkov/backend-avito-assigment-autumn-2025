package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/dto"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service"
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
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "error": map[string]string{
                "code":    "BAD_REQUEST",
                "message": "invalid request body",
            },
        })
        return
    }

    pr, err := h.prService.Create(r.Context(), body.PullRequestID, body.AuthorID, body.PullRequestName)
    if err != nil {

        switch err {
        case service.ErrNotFound:
            w.WriteHeader(http.StatusNotFound)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "error": map[string]string{
                    "code":    service.ErrNotFound.Error(),
                    "message": "resource not found",
                },
            })
            return

        case service.ErrPullRequestExists:
            w.WriteHeader(http.StatusConflict)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "error": map[string]string{
                    "code":    service.ErrPullRequestExists.Error(),
                    "message": "PR id already exists",
                },
            })
            return
        }
        w.WriteHeader(http.StatusInternalServerError)
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
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "error": map[string]string{
                "code":    "BAD_REQUEST",
                "message": "invalid request body",
            },
        })
        return
    }

    pr, replacedBy, err := h.prService.ReassignReviewer(r.Context(), body.PullRequestID, body.OldReviewerID)
    if err != nil {
        switch err {

        case service.ErrNotFound:
            w.WriteHeader(http.StatusNotFound)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "error": map[string]string{
                    "code":    service.ErrNotFound.Error(),
                    "message": "resource not found",
                },
            })
            return

        case service.ErrPullRequestMerged:
            w.WriteHeader(http.StatusConflict)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "error": map[string]string{
                    "code":    service.ErrPullRequestMerged.Error(),
                    "message": "cannot reassign on merged PR",
                },
            })
            return

        case service.ErrNotAssigned:
            w.WriteHeader(http.StatusConflict)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "error": map[string]string{
                    "code":    service.ErrNotAssigned.Error(),
                    "message": "reviewer is not assigned to this PR",
                },
            })
            return

        case service.ErrNoCandidate:
            w.WriteHeader(http.StatusConflict)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "error": map[string]string{
                    "code":    service.ErrNoCandidate.Error(),
                    "message": "no active replacement candidate in team",
                },
            })
            return
        }

        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "pr":          pr,
        "replaced_by": replacedBy,
    })
}

