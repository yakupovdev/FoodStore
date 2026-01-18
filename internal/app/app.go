package app

import "github.com/yakupovdev/FoodStore/internal/service"

const (
	addr = ":9000"
)

type App struct {
	Server *service.Server
}

func NewApp(server *service.Server) *App {
	return &App{
		Server: server,
	}
}
