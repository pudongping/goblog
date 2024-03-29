package models

import (
	"time"

	"github.com/pudongping/goblog/pkg/types"
)

type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;not null"`

	CreatedAt time.Time `gorm:"column:created_at;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;index"`

	// 支持 gorm 软删除
	// DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" sql:"index"`
}

// GetStringID 获取 ID 的字符串格式
func (a BaseModel) GetStringID() string {
	return types.Uint64ToString(a.ID)
}
