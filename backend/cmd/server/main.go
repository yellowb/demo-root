package main

import (
	"context"
	"flag"
	"log"

	"agent-harness-demo/backend/internal/config"
	"agent-harness-demo/backend/internal/httpapi"
	"agent-harness-demo/backend/internal/store"
	"agent-harness-demo/backend/internal/todos"
	"github.com/gin-gonic/gin"
)

func main() {
	var bootstrapOnly bool
	flag.BoolVar(&bootstrapOnly, "bootstrap-only", false, "initialize the database and exit")
	flag.Parse()

	cfg := config.Load()
	gin.SetMode(gin.ReleaseMode)

	db, err := store.Bootstrap(context.Background(), cfg.DBPath)
	if err != nil {
		log.Fatalf("bootstrap backend: %v", err)
	}
	defer db.Close()

	if bootstrapOnly {
		log.Printf("database ready at %s", cfg.DBPath)
		return
	}

	repo := todos.NewRepository(db)
	service := todos.NewService(repo)
	router := httpapi.NewRouter(service)

	log.Printf("backend listening on %s", cfg.Addr)
	if err := router.Run(cfg.Addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
