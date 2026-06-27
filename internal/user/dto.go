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

type MeResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Wallet  string `json:"wallet"`
	Balance string `json:"balance"`
}

func (user *User) Convert() *MeResponse {

	return &MeResponse{
		ID:      user.ID.String(),
		Email:   user.Email,
		Wallet:  user.WithdrawAddress,
		Balance: user.Balance.String(),
	}
}
