package _interface

import "Authentication_Service/internal/model"

type IUserRepository interface {
	GetById(id string) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) (*model.User, error)
}
