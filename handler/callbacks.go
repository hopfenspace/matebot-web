package handler

import (
	"github.com/hopfenspace/matebot-web/conf"
	"github.com/hopfenspace/matebot-web/sdk"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

type Callback struct {
	DB     *gorm.DB
	Config *conf.Config
	Client sdk.SDK
	WP     worker.Pool
}

func (cb *Callback) Create(c echo.Context) error {
	return nil
}

func (cb *Callback) Update(c echo.Context) error {
	return nil
}

func (cb *Callback) Delete(c echo.Context) error {
	return nil
}
