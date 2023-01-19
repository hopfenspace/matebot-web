package models

import (
	"github.com/myOmikron/echotools/utilitymodels"
	"time"
)

type CoreUser struct {
	ID        uint64                  `gorm:"primarykey" json:"-"`
	UserID    uint64                  `gorm:"uniqueIndex:index_unique_core_user" json:"user_id"`
	User      utilitymodels.LocalUser `json:"-"`
	CoreID    uint64                  `gorm:"uniqueIndex:index_unique_core_user" json:"core_id"`
	CreatedAt time.Time               `json:"-"`
	UpdatedAt time.Time               `json:"-"`
}
