package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	httpdelivery "github.com/yakupovdev/FoodStore/internal/delivery/http"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/handler"
	"github.com/yakupovdev/FoodStore/internal/infrastructure/email"
	"github.com/yakupovdev/FoodStore/internal/infrastructure/postgres"
	"github.com/yakupovdev/FoodStore/internal/infrastructure/security"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

func main() {
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	database := os.Getenv("DATABASE")
	host := os.Getenv("HOST")
	portStr := os.Getenv("PORT")
	port, _ := strconv.ParseUint(portStr, 10, 16)
	username := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	conn, err := postgres.NewConnection(appCtx, postgres.Config{
		Database: database,
		Host:     host,
		Port:     uint16(port),
		Username: username,
		Password: password,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close(appCtx)

	if err := postgres.InitSchema(appCtx, conn); err != nil {
		panic(err)
	}

	userRepo := postgres.NewUserRepo(conn)
	tokenRepo := postgres.NewTokenRepo(conn)
	recoveryCodeRepo := postgres.NewRecoveryCodeRepo(conn)
	clientRepo := postgres.NewClientRepo(conn)
	orderRepo := postgres.NewOrderRepo(conn)
	productRepo := postgres.NewProductRepo(conn)
	sellerRepo := postgres.NewSellerRepo(conn)

	codeHasher := security.NewSHA256CodeHasher()
	tokenSvc := security.NewJWTService()
	codeGen := security.NewRandomCodeGenerator()
	emailSender := email.NewSMTPSender(
		"foodstorewwgo@gmail.com",
		"dkeywmbvieuiuazj",
		"smtp.gmail.com",
		"587",
	)

	authUsecase, _ := usecase.NewAuthUsecase(userRepo, tokenRepo, tokenSvc)
	recoveryUsecase, _ := usecase.NewRecoveryUsecase(userRepo, recoveryCodeRepo, codeHasher, tokenSvc, codeGen, emailSender)
	clientUsecase, _ := usecase.NewClientUsecase(clientRepo, orderRepo, productRepo, sellerRepo)
	sellerUsecase, _ := usecase.NewSellerUsecase(sellerRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	emailHandler := handler.NewEmailHandler(recoveryUsecase)
	recoveryHandler := handler.NewRecoveryHandler(recoveryUsecase)
	refreshTokenHandler := handler.NewRefreshTokenHandler(authUsecase)
	clientHandler := handler.NewClientHandler(clientUsecase)
	sellerHandler := handler.NewSellerHandler(sellerUsecase)

	r := httpdelivery.SetupRouter(httpdelivery.RouterDeps{
		AuthHandler:         authHandler,
		EmailHandler:        emailHandler,
		RecoveryHandler:     recoveryHandler,
		RefreshTokenHandler: refreshTokenHandler,
		ClientHandler:       clientHandler,
		SellerHandler:       sellerHandler,
		TokenValidator:      authUsecase,
		TokenService:        tokenSvc,
	})

	srv := &http.Server{
		Addr:         ":9000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Starting server on :9000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-appCtx.Done():
				log.Println("background cleanup stopped")
				return
			case <-ticker.C:
				if err := authUsecase.DeleteExpiredTokens(appCtx); err != nil {
					log.Printf("Error deleting expired access tokens: %v", err)
				} else {
					log.Println("Expired access tokens deleted successfully")
				}
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	appCancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("server stopped gracefully")
}
