package models

type Entity struct {
	ID          uint              `json:"id" gorm:"primary_key"`
	Name        string            `json:"name" gorm:"type:varchar(254);unique"`
	Address     string            `json:"address" gorm:"type:text"`
	Users       []User            `gorm:"many2many:user_entities;"`
	AuditFields `gorm:"embedded"` // Embedding AuditFields
}

// TableName specifies the table name for the Entity model.
func (Entity) TableName() string {
	return "entity"
}

type UserEntity struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	UserID      uint              `json:"user_id" gorm:"index"`   // Foreign key for User
	EntityID    uint              `json:"entity_id" gorm:"index"` // Foreign key for Entity
	AuditFields `gorm:"embedded"` // Embedding AuditFields
}
