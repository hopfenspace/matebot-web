package handler

import (
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/hopfenspace/matebot-web/conf"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

type API struct {
	DB         *gorm.DB
	Config     *conf.Config
	WorkerPool worker.Pool
	SDK        *MateBotSDKGo.SDK
}

type GenericResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
