package server

import (
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/hopfenspace/matebot-web/conf"
	"github.com/hopfenspace/matebot-web/handler"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
	"path/filepath"
)

func defineRoutes(e *echo.Echo, db *gorm.DB, config *conf.Config, client MateBotSDKGo.SDK, wp worker.Pool) {
	api := handler.API{
		DB:            db,
		Config:        config,
		WorkerPool:    wp,
		SDK:           client,
		EventChannels: &map[string]chan handler.EventNotification{},
	}

	e.GET("/", handler.Index)

	e.POST("/api/frontend/login", api.Login)
	e.GET("/api/frontend/logout", api.Logout)
	e.POST("/api/frontend/register", api.Register) // for new users
	e.POST("/api/frontend/connect", api.Connect)   // for existing users
	e.GET("/api/frontend/test", api.Test)

	e.GET("/api/websocket", api.WebSocket)

	e.POST("/api/callback", api.Callback)

	e.Static("/static", filepath.Join(config.Server.StaticDir))
}
