package server

import (
	"github.com/hopfenspace/MateBotSDKGo"
	"github.com/hopfenspace/matebot-web/conf"
	"github.com/hopfenspace/matebot-web/handler"
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

func defineRoutes(e *echo.Echo, db *gorm.DB, config *conf.Config, client MateBotSDKGo.SDK, wp worker.Pool) {
	api := handler.NewAPI(db, config, &client, wp)

	e.POST("/api/frontend/login", api.Login)
	e.GET("/api/frontend/logout", api.Logout)
	e.POST("/api/frontend/register", api.Register)      // for new users
	e.POST("/api/frontend/connectAccount", api.Connect) // for existing users
	e.GET("/api/frontend/test", api.Test)
	e.POST("/api/frontend/test", api.Test)

	e.GET("/api/frontend/consumables", api.Consumables)
	e.GET("/api/frontend/applications", api.Applications)
	e.GET("/api/frontend/blame", api.Blame)
	e.GET("/api/frontend/zwegat", api.Zwegat)
	e.POST("/api/frontend/balance", api.Balance) // primarily for querying other users, use /me otherwise

	e.GET("/api/frontend/me", api.Me)
	e.GET("/api/frontend/listUsers", api.ListUsers)
	e.POST("/api/frontend/changeUsername", api.ChangeUsername) // only changes global username, not local username
	e.POST("/api/frontend/startVouching", api.StartVouching)
	e.POST("/api/frontend/stopVouching", api.StopVouching)
	e.POST("/api/frontend/dropPrivileges", api.DropPrivileges)
	e.POST("/api/frontend/confirmAlias", api.ConfirmAlias)
	e.POST("/api/frontend/deleteAlias", api.DeleteAlias)
	e.POST("/api/frontend/deleteLocalAccount", api.DeleteLocalAccount)
	e.POST("/api/frontend/deleteFullAccount", api.DeleteFullAccount)

	e.GET("/api/frontend/getTransactions", api.ListTransactions) // only own transaction history
	e.POST("/api/frontend/sendTransaction", api.SendTransaction)
	e.POST("/api/frontend/consumeTransaction", api.ConsumeTransaction)

	e.GET("/api/frontend/openCommunisms", api.OpenCommunisms)
	e.GET("/api/frontend/allCommunisms", api.AllCommunisms)
	e.POST("/api/frontend/newCommunism", api.NewCommunism)
	e.POST("/api/frontend/closeCommunism", api.CloseCommunism)
	e.POST("/api/frontend/joinCommunism", api.JoinCommunism)
	e.POST("/api/frontend/leaveCommunism", api.LeaveCommunism)
	e.POST("/api/frontend/abortCommunism", api.AbortCommunism)

	e.GET("/api/frontend/openPolls", api.OpenPolls)
	e.GET("/api/frontend/allPolls", api.AllPolls)
	e.POST("/api/frontend/newPoll", api.NewPoll)
	e.POST("/api/frontend/approvePoll", api.ApprovePoll)
	e.POST("/api/frontend/disapprovePoll", api.DisapprovePoll)
	e.POST("/api/frontend/abortPoll", api.AbortPoll)

	e.GET("/api/frontend/openRefunds", api.OpenRefunds)
	e.GET("/api/frontend/allRefunds", api.AllRefunds)
	e.POST("/api/frontend/newRefund", api.NewRefund)
	e.POST("/api/frontend/approveRefund", api.ApproveRefund)
	e.POST("/api/frontend/disapproveRefund", api.DisapproveRefund)
	e.POST("/api/frontend/abortRefund", api.AbortRefund)

	e.GET("/api/websocket", api.WebSocket)

	e.POST("/api/callback", api.Callback)
}
