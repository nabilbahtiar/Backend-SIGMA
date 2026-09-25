package repository

import (
	"server-room-auth/internal/model"
	"server-room-auth/pkg/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByNIK(nik string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("nik = ?", nik).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	return database.DB.Save(user).Error
}
