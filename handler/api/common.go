package api

import (
	"github.com/hopfenspace/matebot-web/conf"
	"github.com/hopfenspace/matebot-web/sdk"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

type Api struct {
	DB         *gorm.DB
	Config     *conf.Config
	WorkerPool worker.Pool
	Client     sdk.SDK
}
