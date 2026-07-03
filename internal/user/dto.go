package user

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Wallet   string `json:"wallet"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateRequest struct {
	Wallet string `json:"wallet" binding:"required"`
}

type Response struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Wallet  string `json:"wallet"`
	Balance string `json:"balance"`
}
