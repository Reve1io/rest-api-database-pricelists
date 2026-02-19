package main

import (
	"context"
	"log"
	"net/http"

	"rest-api-database-pricelists/internal/config"
	"rest-api-database-pricelists/internal/handler"
	"rest-api-database-pricelists/internal/logger"
	"rest-api-database-pricelists/internal/repository"
	"rest-api-database-pricelists/internal/service"
	"rest-api-database-pricelists/pkg/db"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	godotenv.Load()

	cfg := config.Load()

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	logg, _ := logger.New()

	repo := repository.NewProductRepository(pool, logg)
	svc := service.NewSearchService(repo, logg)
	h := handler.NewSearchHandler(svc, logg)

	r := chi.NewRouter()

	r.Post("/api/v1/search", h.Search)

	port := cfg.AppPort

	logg.Info("server starting",
		zap.String("port:", port),
	)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
