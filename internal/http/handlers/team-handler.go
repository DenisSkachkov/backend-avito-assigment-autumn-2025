package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/team"
	"github.com/gorilla/mux"
)

type TeamHandler struct {
	service *team.TeamService
}

func NewTeamHandler(service *team.TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

func(h *TeamHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/team/add", h.CreateTeam).Methods(http.MethodPost)
	r.HandleFunc("/team/get", h.GetTeam).Methods(http.MethodGet)
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var t models.Team
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateTeam(r.Context(), &t)
	if err != nil {
		if errors.Is(err, service.ErrTeamExists) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
				},
			})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"team": created,
	})
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		http.Error(w, "team_name is required", http.StatusBadRequest)
		return
	}

	t, err := h.service.GetTeam(r.Context(), teamName)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "error": map[string]string{
                    "code":    service.ErrNotFound.Error(),
                    "message": "resource not found",
                },
            })
            return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(t)
}