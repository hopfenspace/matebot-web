package server

import (
	"github.com/hopfenspace/matebot-web/conf"
	"github.com/hopfenspace/matebot-web/handler"
	"github.com/hopfenspace/matebot-web/handler/api"
	"github.com/hopfenspace/matebot-web/sdk"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
	"path/filepath"
)

func defineRoutes(e *echo.Echo, db *gorm.DB, config *conf.Config, client sdk.SDK, wp worker.Pool) {
	e.GET("/", handler.Index)

	// API for frontend
	frontend := api.Api{
		DB:         db,
		Config:     config,
		WorkerPool: wp,
		Client:     client,
	}
	e.POST("/api/frontend/login", frontend.Login)
	e.GET("/api/frontend/logout", frontend.Logout)
	e.POST("/api/frontend/register", frontend.Register)
	e.POST("/api/frontend/test", frontend.Test)

	// Callbacks for MateBot Core
	cb := handler.Callback{
		Config: config,
		DB:     db,
		Client: client,
		WP:     wp,
	}
	e.GET("/core/create/:model/:id", cb.Create)
	e.GET("/core/update/:model/:id", cb.Update)
	e.GET("/core/delete/:model/:id", cb.Delete)

	e.Static("/static", filepath.Join(config.Generic.StaticDir))
}
