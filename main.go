package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"user-service/docs"
	"user-service/handlers"
)

func router(h *handlers.Handler) *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.Title = "Users service"
	docs.SwaggerInfo.Description = "This is a sample users service"
	docs.SwaggerInfo.Version = "1.0"

	r.GET("/users/:uuid", h.GetUser())
	r.POST("/users", h.CreateUser())
	r.PUT("/users/:uuid", h.ChangeUser())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func main() {
	env := LoadEnv()

	db, err := sql.Open("postgres", env.Db.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("failed connect to db", err)
	}
	handler := handlers.NewHandler(db)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", env.App.Port),
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
