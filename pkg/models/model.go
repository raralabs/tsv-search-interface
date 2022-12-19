package models

import (
	"github.com/aymericbeaumet/go-tsvector"
	"gorm.io/datatypes"
)

type SearchIndex struct {
	ID          string            `json:"id"`
	TableInfo   string            `json:"table_info"`
	ActionInfo  datatypes.JSON    `gorm:"type:jsonb" json:"action_info"`
	SearchField datatypes.JSON    `gorm:"type:jsonb" json:"search_field"`
	TsvText     tsvector.TSVector `gorm:"not null"`
}

type ResponseSearchIndex struct {
	ID         string         `json:"id"`
	TableInfo  string         `json:"table_info"`
	ActionInfo datatypes.JSON `json:"action_info"`
}
