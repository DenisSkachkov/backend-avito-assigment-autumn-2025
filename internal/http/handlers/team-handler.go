package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/team"
	"github.com/gorilla/mux"
)

type TeamHandler struct {
	service team.TeamService
}

func NewTeamHandler(service team.TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

func(h *TeamHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/team/add", h.CreateTeam).Methods(http.MethodPost)
	r.HandleFunc("/team/get", h.GetTeam).Methods(http.MethodGet)
}

func(h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team

    if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    created, err := h.service.CreateTeam(r.Context(), &team)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
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

    team, err := h.service.GetTeam(r.Context(), teamName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(team)
}