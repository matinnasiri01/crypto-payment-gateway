package user

import "context"

type Repository interface {
	Create(cnx context.Context, user *User) error
	Update(cnx context.Context, user *User) error
	GetByID(cnx context.Context, id string) (*User, error)
	GetByEmail(cnx context.Context, email string) (*User, error)
}
