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
	"github.com/yakupovdev/FoodStore/internal/infrastructure/postgres/impl"
	"github.com/yakupovdev/FoodStore/internal/infrastructure/postgres/initializationdb"

	httpdelivery "github.com/yakupovdev/FoodStore/internal/delivery/http"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/handler"
	"github.com/yakupovdev/FoodStore/internal/infrastructure/email"
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

	conn, err := initializationdb.NewConnection(appCtx, initializationdb.Config{
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

	if err := initializationdb.InitSchema(appCtx, conn); err != nil {
		panic(err)
	}

	// Repository
	userRepo := impl.NewUserRepo(conn)
	tokenRepo := impl.NewTokenRepo(conn)
	recoveryCodeRepo := impl.NewRecoveryCodeRepo(conn)
	clientRepo := impl.NewClientRepo(conn)
	orderRepo := impl.NewOrderRepo(conn)
	productRepo := impl.NewProductRepo(conn)
	sellerRepo := impl.NewSellerRepo(conn)
	transactionRepo := impl.NewTransactionRepository(conn)

	// Services
	codeHasher := security.NewSHA256CodeHasher()
	tokenSvc := security.NewJWTService()
	codeGen := security.NewRandomCodeGenerator()
	emailSender := email.NewSMTPSender(
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PORT"),
	)

	// Usecases
	authUsecase, _ := usecase.NewAuthUsecase(userRepo, tokenRepo, tokenSvc)
	recoveryUsecase, _ := usecase.NewRecoveryUsecase(userRepo, recoveryCodeRepo, codeHasher, tokenSvc, codeGen, emailSender)
	clientUsecase, _ := usecase.NewClientUsecase(clientRepo, orderRepo, productRepo, sellerRepo, transactionRepo)
	sellerUsecase, _ := usecase.NewSellerUsecase(sellerRepo, productRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	emailHandler := handler.NewEmailHandler(recoveryUsecase)
	recoveryHandler := handler.NewRecoveryHandler(recoveryUsecase)
	refreshTokenHandler := handler.NewRefreshTokenHandler(authUsecase)
	clientHandler := handler.NewClientHandler(clientUsecase)
	sellerHandler := handler.NewSellerHandler(sellerUsecase)

	// Router
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

	// Server
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
				if err := authUsecase.DeleteExpiredTokens(appCtx); err != nil {
					log.Printf("Error deleting expired access tokens: %v", err)
				} else {
					log.Println("Expired access tokens deleted successfully")
				}
			}
		}
	}()

	// Gracefully shutdown
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
