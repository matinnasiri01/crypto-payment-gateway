package user

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Wallet   string `json:"wallet"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MeResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Wallet  string `json:"wallet"`
	Balance string `json:"balance"`
}
