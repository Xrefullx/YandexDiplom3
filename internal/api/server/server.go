package server

import (
	"context"
	"errors"
	"github.com/Xrefullx/YandexDiplom3/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func InitServer(r *gin.Engine, cfg models.Config) {
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ATTENTION!!!!: %s\n\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("COMING SOON WE BACK", err)
	}
	log.Println("I`M GO SLEEP")
}
