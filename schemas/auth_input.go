package schemas

type AuthInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserStruct struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"` // Omit from JSON responses
	Entities  string `json:"enitities"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
	UpdatedAt string `json:"updated_at"`
	UpdatedBy string `json:"updated_by"`
	Updator   string `json:"updator"`
	IsActive  string `json:"is_active"`
}
