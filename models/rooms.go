package models

import "time"

type RoomCategory struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Name     string  `json:"name" gorm:"type:varchar(100);unique;not null"`
	Price    float64 `json:"price" gorm:"type:decimal(10,2);not null"` // Default price for category (can be overridden by specific rooms)
	IsActive bool    `json:"is_active" gorm:"default:true"`
	AuditFields
}

func (RoomCategory) TableName() string {
	return "room_categories"
}

type Room struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	RoomNumber     string       `json:"room_number" gorm:"type:varchar(50);unique;not null"`
	RoomCategoryID uint         `json:"room_category_id" gorm:"not null"`
	RoomCategory   RoomCategory `gorm:"foreignKey:RoomCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Price          float64      `json:"price" gorm:"type:decimal(10,2);not null"` // Price per day for the room
	IsActive       bool         `json:"is_active" gorm:"default:true"`
	AuditFields
}

func (Room) TableName() string {
	return "rooms"
}

type RoomOccupancy struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	RoomID        uint       `json:"room_id" gorm:"not null"`
	Room          Room       `gorm:"foreignKey:RoomID"`
	PatientID     *uint      `json:"patient_id" gorm:"index;default:null"`
	Patient       *Patient   `gorm:"foreignKey:PatientID"`
	AdmitDate     time.Time  `json:"admit_date" gorm:"not null"`
	DischargeDate *time.Time `json:"discharge_date" gorm:"default:null"`
	PricePerDay   float64    `json:"price_per_day" gorm:"type:decimal(10,2);not null"` // The price the patient will be charged per day
	IsActive      bool       `json:"is_active" gorm:"default:true"`                    // True if patient is still admitted
	AuditFields
}

func (RoomOccupancy) TableName() string {
	return "room_occupancy"
}
