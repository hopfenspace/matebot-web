package models

import (
	"github.com/myOmikron/echotools/utilitymodels"
	"time"
)

type CoreUser struct {
	ID        uint                    `gorm:"primarykey" json:"-"`
	UserID    uint                    `gorm:"unique" json:"user_id"`
	User      utilitymodels.LocalUser `json:"-"`
	MateBotID uint                    `gorm:"unique" json:"core_id"`
	CreatedAt time.Time               `json:"-"`
	UpdatedAt time.Time               `json:"-"`
}
