package models

import "encoding/json"

// Import this for JSONB support

type MedicineCategory struct {
	ID               uint              `json:"id" gorm:"primaryKey"`
	Name             string            `json:"name" gorm:"type:varchar(100);unique"`
	NameTranslations json.RawMessage `gorm:"type:jsonb" json:"name_translations"` 
	Description      string            `json:"description"`
	EntityID         uint              `json:"entity_id"` // Foreign key to Entity
	Entity           Entity            `json:"entity" gorm:"foreignKey:EntityID"`
	Medicines        []Medicine        `json:"medicines" gorm:"foreignKey:CategoryID"`
	IsActive         bool              `json:"is_active" gorm:"default:true"`
	AuditFields      `gorm:"embedded"` // Embedding AuditFields
}

func (MedicineCategory) TableName() string {
	return "medicine_category"
}

type Medicine struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	Name        string            `json:"name" gorm:"type:varchar(100);unique"`
	NameTranslations json.RawMessage `gorm:"type:jsonb" json:"name_translations"` 
	CategoryID  uint              `json:"category_id"` // Foreign key to MedicineCategory
	Category    MedicineCategory  `json:"category" gorm:"foreignKey:CategoryID"`
	EntityID    uint              `json:"entity_id"` // Foreign key to Entity
	Entity      Entity            `json:"entity" gorm:"foreignKey:EntityID"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Quantity       int               `json:"quantity"`
	IsActive    bool              `json:"is_active" gorm:"default:true"`
	AuditFields `gorm:"embedded"` // Embedding AuditFields
}

func (Medicine) TableName() string {
	return "medicine"
}
