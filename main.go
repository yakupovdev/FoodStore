package main

import (
	"context"
	"log"
	"strconv"

	"os"

	"github.com/joho/godotenv"
	"github.com/yakupovdev/FoodStore/internal/app"
	"github.com/yakupovdev/FoodStore/internal/repository"
	"github.com/yakupovdev/FoodStore/internal/service"
	"github.com/yakupovdev/FoodStore/internal/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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

	pgStorage := repository.NewPostgres(conn)

	srv := service.NewServer(":9000", service.SetupRouter(pgStorage))
	application := app.NewApp(srv)

	log.Println("Starting server on :9000")
	if err := application.Server.Run(); err != nil {
		panic(err)
	}

}
