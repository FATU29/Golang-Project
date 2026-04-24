package model

import "time"

type User struct {
	Id        string    `json:"id" gorm:"column:id;primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Firstname *string   `json:"firstname"`
	Lastname  *string   `json:"lastname"`
	Password  string    `json:"password"`
	Avatar     *string   `json:"avatar"`
	CoverImage *string   `json:"cover_image"`
	IsActive   int       `json:"is_active"`
}

func (User) TableName() string {
	return "users"
}
