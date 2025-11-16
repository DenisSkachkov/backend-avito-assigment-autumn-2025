package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/http"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/http/handlers"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/repository/postgres"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/pullrequest"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/team"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/user"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		log.Fatal("Failed to connect postgres:", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepo(db)
	teamRepo := postgres.NewTeamRepo(db)
	prRepo := postgres.NewPRRepo(db)

	userService := user.NewUserService(userRepo)
	teamService := team.NewTeamService(teamRepo,userRepo)
	prService := pullrequest.NewPullRequestService(userRepo,prRepo,teamRepo)

	r := mux.NewRouter()

	handlers.NewPullRequestHandler(*prService).RegisterRoutes(r)
	handlers.NewTeamHandler(*teamService).RegisterRoutes(r)
	handlers.NewUserHandler(*userService, *prService).RegisterRoutes(r)

	server := http.New(r, "8080")

	go func() {
        if err := server.Start(); err != nil {
            log.Fatal(err)
        }
    }()

	stop := make(chan os.Signal, 1)
    <-stop

    ctxTimeout, cancel := context.WithTimeout(ctx, 5* time.Second)
    defer cancel()

    server.Stop(ctxTimeout)
}