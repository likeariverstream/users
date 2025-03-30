package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	users "user-service"
)

var storage = make(map[string]string)

func router(h *users.Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/users/:uuid", h.GetUser())

	r.POST("/users", h.CreateUser())

	r.PUT("/users/:uuid", h.ChangeUser())

	return r
}

func main() {
	handler := users.NewHandler(storage)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router(handler),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("server exiting")

}
