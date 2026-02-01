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
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	// DB
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading ..env file")
	}
	Database := os.Getenv("DATABASE")
	Host := os.Getenv("HOST")
	PortBefore := os.Getenv("PORT")
	Port, _ := strconv.ParseUint(PortBefore, 10, 16)
	Username := os.Getenv("USER")
	Password := os.Getenv("PASSWORD")

	conn, err := storage.NewPostgresDB(appCtx, storage.Config{
		Database: Database,
		Host:     Host,
		Port:     uint16(Port),
		Username: Username,
		Password: Password,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close(appCtx)

	if err := storage.InitSchema(appCtx, conn); err != nil {
		panic(err)
	}

	// Repository
	userRepo := repository.NewPostgres(conn)
	clientRepo := repository.NewOrdersRepo(conn)
	sellerRepo := repository.NewSellerRepository(conn)

	// Usecases
	authUsecase, _ := usecase.NewAuthUsecase(userRepo)
	emailUsecase, _ := usecase.NewEmailUsecase(userRepo)
	recoveryUsecase, _ := usecase.NewRecoveryUsecase(userRepo)
	refreshAccessTokenUsecase, _ := usecase.NewRefreshAccessTokenUsecase(userRepo)
	chechTokenIsValidUsecase, _ := usecase.NewCheckTokenIsValidUsecase(userRepo)
	clientUsecase, _ := usecase.NewClientUsecase(clientRepo)
	sellerUsecase, _ := usecase.NewSellerUsecase(sellerRepo)

	// Controllers
	authController := controller.NewAuthController(authUsecase)
	emailController := controller.NewEmailController(emailUsecase)
	recoveryController := controller.NewRecoveryController(recoveryUsecase)
	refreshAccessTokenController := controller.NewRefreshAccessTokenController(refreshAccessTokenUsecase)
	clientController := controller.NewClientController(clientUsecase)
	sellerController := controller.NewSellerController(sellerUsecase)

	// Router
	r := router.SetupRouter(router.Deps{
		AuthController:               authController,
		EmailController:              emailController,
		RecoveryController:           recoveryController,
		RefreshAccessTokenController: refreshAccessTokenController,
		CheckTokenIsValidUsecase:     chechTokenIsValidUsecase,
		ClientController:             clientController,
		SellerController:             sellerController,
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

	// Background task
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-appCtx.Done():
				log.Println("background cleanup stopped")
				return
			case <-ticker.C:
				if err := userRepo.DeleteExpiredAccessTokens(); err != nil {
					log.Printf("Error deleting expired access tokens: %v", err)
				} else {
					log.Println("Expired access tokens deleted successfully")
				}
			}
		}
	}()

	// GRACEFUL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	// остановить фоновые задачи
	appCancel()

	// отдельный контекст только для Shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("server stopped gracefully")
}
