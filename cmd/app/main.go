package main

import (
	"context"
	"log"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/yakupovdev/FoodStore/internal/app"
	"github.com/yakupovdev/FoodStore/internal/service"
	"github.com/yakupovdev/FoodStore/internal/storage"
	"github.com/yakupovdev/FoodStore/pkg"
	"os"
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

	db, err := pkg.NewPostgresDB(ctx, pkg.Config{
		Database: Database,
		HOST:     Host,
		Port:     uint16(Port),
		Username: Username,
		Password: Password,
	})
	if err != nil {
		panic(err)
	}

	defer db.Conn.Close(ctx)
	pgStorage := storage.NewPostgres(db.Conn)

	srv := service.NewServer(":9000", service.SetupRouter(pgStorage))
	application := app.NewApp(srv)

	log.Println("Starting server on :9000")
	if err := application.Server.Run(); err != nil {
		panic(err)
	}

}
