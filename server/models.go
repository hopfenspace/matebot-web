package server

import (
	"github.com/myOmikron/echotools/utilitymodels"
	"time"
)

type BotUser struct {
	ID        uint                    `gorm:"primarykey" json:"-"`
	UserID    uint                    `json:"user_id"`
	User      utilitymodels.LocalUser `json:"-"`
	MateBotID uint                    `json:"bot_id"`
	CreatedAt time.Time               `json:"-"`
	UpdatedAt time.Time               `json:"-"`
}
