package user

type Service struct {
	repo *Repository
}

func New(r *Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (ser Service) SignUp(req *SignupRequest) error {
	return nil
}

func (ser Service) LogIn(req *LoginRequest) error {
	return nil
}

func (ser Service) Update(req *LoginRequest) error {
	return nil
}

func (ser Service) GetUserByID(ID string) (*MeResponse, error) {
	return nil, nil
}
