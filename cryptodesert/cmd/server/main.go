package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crypto-desert/internal/api"
	"crypto-desert/internal/store"
)

func main() {
	// ── Configuração via variáveis de ambiente ────────────────────────────────
	port := envOr("PORT", "8080")
	webDir := envOr("WEB_DIR", "./web")

	// ── Stores (in-memory) ────────────────────────────────────────────────────
	chars := store.NewCharacterStore()
	battles := store.NewBattleStore()
	runners := store.NewRunnerStore()
	inventories := store.NewInventoryStore()

	// ── Serviços ──────────────────────────────────────────────────────────────
	cryptoSvc := api.NewCryptoService()

	// ── Handler e Rotas ───────────────────────────────────────────────────────
	handler := api.NewHandler(chars, battles, runners, inventories, cryptoSvc)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, handler, webDir)

	// ── Servidor HTTP ─────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Inicia servidor em goroutine para poder escutar sinais
	go func() {
		fmt.Printf("\n")
		fmt.Printf("  ██████╗██████╗ ██╗   ██╗██████╗ ████████╗ ██████╗\n")
		fmt.Printf(" ██╔════╝██╔══██╗╚██╗ ██╔╝██╔══██╗╚══██╔══╝██╔═══██╗\n")
		fmt.Printf(" ██║     ██████╔╝ ╚████╔╝ ██████╔╝   ██║   ██║   ██║\n")
		fmt.Printf(" ██║     ██╔══██╗  ╚██╔╝  ██╔═══╝    ██║   ██║   ██║\n")
		fmt.Printf(" ╚██████╗██║  ██║   ██║   ██║        ██║   ╚██████╔╝\n")
		fmt.Printf("  ╚═════╝╚═╝  ╚═╝   ╚═╝   ╚═╝        ╚═╝    ╚═════╝\n")
		fmt.Printf("\n")
		fmt.Printf("  ██████╗ ███████╗███████╗███████╗██████╗ ████████╗\n")
		fmt.Printf("  ██╔══██╗██╔════╝██╔════╝██╔════╝██╔══██╗╚══██╔══╝\n")
		fmt.Printf("  ██║  ██║█████╗  ███████╗█████╗  ██████╔╝   ██║\n")
		fmt.Printf("  ██║  ██║██╔══╝  ╚════██║██╔══╝  ██╔══██╗   ██║\n")
		fmt.Printf("  ██████╔╝███████╗███████║███████╗██║  ██║   ██║\n")
		fmt.Printf("  ╚═════╝ ╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝   ╚═╝\n")
		fmt.Printf("\n")
		fmt.Printf("  RPG · 2087 · Deserto Digital\n")
		fmt.Printf("  ──────────────────────────────────────────────────\n")
		fmt.Printf("  Servidor:  http://localhost:%s\n", port)
		fmt.Printf("  Frontend:  http://localhost:%s\n", port)
		fmt.Printf("  API Base:  http://localhost:%s/api\n", port)
		fmt.Printf("  Web Dir:   %s\n", webDir)
		fmt.Printf("  ──────────────────────────────────────────────────\n\n")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[server] fatal: %v", err)
		}
	}()

	// ── Graceful Shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[server] shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[server] forced shutdown: %v", err)
	}
	log.Println("[server] stopped.")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
