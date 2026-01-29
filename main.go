package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"os"

	"github.com/joho/godotenv"
	"github.com/yakupovdev/FoodStore/internal/controller"
	"github.com/yakupovdev/FoodStore/internal/repository"
	"github.com/yakupovdev/FoodStore/internal/router"
	"github.com/yakupovdev/FoodStore/internal/storage"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

func main() {
	// DB
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading ..env file")
	}
	Database := os.Getenv("DATABASE")
	Host := os.Getenv("HOST")
	PortBefore := os.Getenv("PORT")
	Port, _ := strconv.ParseUint(PortBefore, 10, 16)
	Username := os.Getenv("USER")
	Password := os.Getenv("PASSWORD")

	ctx := context.Background()

	conn, err := storage.NewPostgresDB(ctx, storage.Config{
		Database: Database,
		Host:     Host,
		Port:     uint16(Port),
		Username: Username,
		Password: Password,
	})
	if err != nil {
		panic(err)
	}

	if err := storage.InitSchema(ctx, conn); err != nil {
		panic(err)
	}

	defer conn.Close(ctx)

	// Repository
	repo := repository.NewPostgres(conn)

	// Usecases
	authUsecase, _ := usecase.NewAuthUsecase(repo)
	emailUsecase, _ := usecase.NewEmailUsecase(repo)
	recoveryUsecase, _ := usecase.NewRecoveryUsecase(repo)
	refreshAccessTokenUsecase, _ := usecase.NewRefreshAccessTokenUsecase(repo)
	chechTokenIsValidUsecase, _ := usecase.NewCheckTokenIsValidUsecase(repo)

	// Controllers
	authController := controller.NewAuthController(authUsecase)
	emailController := controller.NewEmailController(emailUsecase)
	recoveryController := controller.NewRecoveryController(recoveryUsecase)
	refreshAccessTokenController := controller.NewRefreshAccessTokenController(refreshAccessTokenUsecase)

	// Router
	r := router.SetupRouter(router.Deps{
		AuthController:               authController,
		EmailController:              emailController,
		RecoveryController:           recoveryController,
		RefreshAccessTokenController: refreshAccessTokenController,
		CheckTokenIsValidUsecase:     chechTokenIsValidUsecase,
	})

	// HTTP Server
	srv := &http.Server{
		Addr:         ":9000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run Server
	go func() {
		log.Println("Starting server on :9000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// GRACEFUL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("server stopped gracefully")
}
