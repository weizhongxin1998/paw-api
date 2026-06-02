package main

import (
	"context"
	"log"

	"paw-api/database"
	"paw-api/pkg/snowflake"
)

type App struct {
	ctx       context.Context
	snowflake *snowflake.Generator
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.snowflake = snowflake.New()

	db, err := database.Init()
	if err != nil {
		log.Fatal("failed to init database:", err)
	}
	_ = db
}
