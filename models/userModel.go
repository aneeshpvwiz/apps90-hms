package models

type User struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Email    string    `json:"email" gorm:"type:varchar(254);unique"`
	Password string    `json:"password" gorm:"type:varchar(254)"`
	Entities []*Entity `gorm:"many2many:user_entities;"`
	AuditFields
}

// TableName specifies the table name for the User model.
func (User) TableName() string {
	return "auth_user"
}
