package schemas

type EntityInput struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type UserEntityInput struct {
	UserID   uint `json:"user_id" binding:"required"`
	EntityID uint `json:"entity_id" binding:"required"`
}
