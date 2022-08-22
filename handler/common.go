package api

import (
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/hopfenspace/matebot-web/conf"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

type Api struct {
	DB         *gorm.DB
	Config     *conf.Config
	WorkerPool worker.Pool
	SDK        *MateBotSDKGo.SDK
}
