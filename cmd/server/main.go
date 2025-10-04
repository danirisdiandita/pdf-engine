package main

import (
	"log"

	"github.com/danirisdiandita/pdf-engine/internal/config"
	"github.com/danirisdiandita/pdf-engine/internal/router"
)

func main() {
	cfg := config.Load()

	r := router.SetupRouter(cfg)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
