package user

import "context"

type Repository interface {
	Create(cnx context.Context, user *User) error
	Update(cnx context.Context, user *User) error
	GetUserByID(cnx context.Context, id string) (*User, error)
}
