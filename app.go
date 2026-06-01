package main

import (
	"context"
	"paw-api/database"
	"paw-api/handlers"
)

type App struct {
	ctx                context.Context
	ProjectHandler     *handlers.ProjectHandler
	CollectionHandler  *handlers.CollectionHandler
	RequestHandler     *handlers.RequestHandler
	EnvironmentHandler *handlers.EnvironmentHandler
	HistoryHandler     *handlers.HistoryHandler
	WebSocketHandler   *handlers.WebSocketHandler
}

func NewApp() *App {
	return &App{
		ProjectHandler:     handlers.NewProjectHandler(),
		CollectionHandler:  handlers.NewCollectionHandler(),
		RequestHandler:     handlers.NewRequestHandler(),
		EnvironmentHandler: handlers.NewEnvironmentHandler(),
		HistoryHandler:     handlers.NewHistoryHandler(),
		WebSocketHandler:   handlers.NewWebSocketHandler(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.WebSocketHandler.SetContext(ctx)
	if err := database.InitDB(); err != nil {
		panic("failed to initialize database: " + err.Error())
	}
}

func (a *App) shutdown(ctx context.Context) {
	database.CloseDB()
}
