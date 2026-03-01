package repository

import (
	"errors"

	"Authentication_Service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	Db *gorm.DB
}

func (u UserRepository) GetById(id string) (*model.User, error) {
	var user model.User
	res := u.Db.Where("id = ?", id).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (u UserRepository) Create(user *model.User) (*model.User, error) {
	if user.Id == "" {
		user.Id = uuid.New().String()
	}
	res := u.Db.Create(user)

	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (u UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	res := u.Db.Where("email = ? ", email).First(&user)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (u UserRepository) Update(user *model.User) (*model.User, error) {
	updates := map[string]interface{}{
		"email":      user.Email,
		"firstname":  user.Firstname,
		"lastname":   user.Lastname,
		"password":   user.Password,
		"avatar":     user.Avatar,
		"is_active":  user.IsActive,
		"updated_at": user.UpdatedAt,
	}
	var res *gorm.DB
	if user.Id != "" {
		res = u.Db.Model(user).Where("id = ?", user.Id).Updates(updates)
	} else {
		res = u.Db.Model(&model.User{}).Where("email = ?", user.Email).Updates(updates)
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
