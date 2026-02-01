package _interface

import "Authentication_Service/internal/model"

type IUserRepository interface {
	getById(id string) (*model.User, error)
	create(user *model.User) (*model.User, error)
}
