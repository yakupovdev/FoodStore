package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/middleware"
	logger "github.com/yakupovdev/FoodStore/internal/domain/logger"
	"github.com/yakupovdev/FoodStore/internal/infrastructure/postgres/impl"
	"github.com/yakupovdev/FoodStore/internal/infrastructure/postgres/initialization"
	"golang.org/x/time/rate"

	httpdelivery "github.com/yakupovdev/FoodStore/internal/delivery/http"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/handler"
	"github.com/yakupovdev/FoodStore/internal/infrastructure/email"
	"github.com/yakupovdev/FoodStore/internal/infrastructure/security"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

func main() {
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()
	limiter := middleware.NewIPLimiter(appCtx, rate.Every(100*time.Millisecond), 20)
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	var pgConfig initialization.Config
	err := envconfig.Process("foodstore", &pgConfig)
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	conn, err := initialization.NewConnection(appCtx, initialization.Config{
		Database: pgConfig.Database,
		Host:     pgConfig.Host,
		Port:     pgConfig.Port,
		User:     pgConfig.User,
		Password: pgConfig.Password,
	})

	if err != nil {
		panic(err)
	}
	defer conn.Close(appCtx)

	if err := initialization.InitSchema(appCtx, conn); err != nil {
		panic(err)
	}

	logger, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
		os.Exit(1)
	}
	defer logger.Close()
	logger.Info("starting FoodStore API server...")

	var smtpConfig email.SMTPSender
	if err := envconfig.Process("email", &smtpConfig); err != nil {
		log.Fatalf("Failed to load email configuration: %v", err)
		os.Exit(1)
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
	moderatorRepo := impl.NewModeratorRepo(conn)
	loggerRepo := impl.NewLogsRepository(conn)
	adminRepo := impl.NewAdminRepository(conn)

	// Services
	codeHasher := security.NewSHA256CodeHasher()
	tokenSvc := security.NewJWTService()
	codeGen := security.NewRandomCodeGenerator()
	emailSender := email.NewSMTPSender(
		smtpConfig.From,
		smtpConfig.Password,
		smtpConfig.Host,
		smtpConfig.Port,
	)
	checkerAdminKey := security.NewCheckerAdminKey(os.Getenv("SECRET_KEY"))

	// Usecases
	authUsecase, _ := usecase.NewAuthUsecase(userRepo, tokenRepo, tokenSvc)
	recoveryUsecase, _ := usecase.NewRecoveryUsecase(userRepo, recoveryCodeRepo, codeHasher, tokenSvc, codeGen, emailSender)
	sellerUsecase, _ := usecase.NewSellerUsecase(sellerRepo, productRepo, moderatorRepo)
	moderatorUsecase, _ := usecase.NewModeratorUsecase(moderatorRepo, productRepo, sellerRepo, emailSender)
	clientUsecase, _ := usecase.NewClientUsecase(clientRepo, orderRepo, productRepo, sellerRepo, transactionRepo, loggerRepo)
	adminUsecase, _ := usecase.NewAdminUsecase(userRepo, clientRepo, adminRepo, loggerRepo, checkerAdminKey)

	// Handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	emailHandler := handler.NewEmailHandler(recoveryUsecase)
	recoveryHandler := handler.NewRecoveryHandler(recoveryUsecase)
	refreshTokenHandler := handler.NewRefreshTokenHandler(authUsecase)
	clientHandler := handler.NewClientHandler(clientUsecase)
	sellerHandler := handler.NewSellerHandler(sellerUsecase)
	moderatorHandler := handler.NewModeratorHandler(moderatorUsecase)
	adminHandler := handler.NewAdminHandler(adminUsecase)

	// Router
	deps := httpdelivery.NewRouterDeps(authHandler, emailHandler, recoveryHandler, refreshTokenHandler, clientHandler, sellerHandler, moderatorHandler, adminHandler, authUsecase, tokenSvc, logger, limiter)
	r := httpdelivery.SetupRouter(deps)

	// Server
	srv := &http.Server{
		Addr:    ":9000",
		Handler: r,

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
				if err := adminUsecase.DeleteExpiredSubscription(appCtx); err != nil {
					log.Printf("Error deleting expired subscriptions: %v", err)
				} else {
					log.Println("Expired subscriptions deleted successfully")
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
