package models

type MedicineCategory struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	Name        string            `json:"name" gorm:"type:varchar(100);unique"`
	Description string            `json:"description"`
	EntityID    uint              `json:"entity_id"` // Foreign key to Entity
	Entity      Entity            `json:"entity" gorm:"foreignKey:EntityID"`
	Medicines   []Medicine        `json:"medicines" gorm:"foreignKey:CategoryID"`
	AuditFields `gorm:"embedded"` // Embedding AuditFields
}

func (MedicineCategory) TableName() string {
	return "medicine_category"
}

type Medicine struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	Name        string            `json:"name" gorm:"type:varchar(100);unique"`
	CategoryID  uint              `json:"category_id"` // Foreign key to MedicineCategory
	Category    MedicineCategory  `json:"category" gorm:"foreignKey:CategoryID"`
	EntityID    uint              `json:"entity_id"` // Foreign key to Entity
	Entity      Entity            `json:"entity" gorm:"foreignKey:EntityID"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Stock       int               `json:"stock"`
	AuditFields `gorm:"embedded"` // Embedding AuditFields
}

func (Medicine) TableName() string {
	return "medicine"
}
