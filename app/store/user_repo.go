package store

import (
	"goals/app/models"
)

// User model's repository
type UserRepo struct {
	store *Store
}


func (r *UserRepo) Create(u *models.User) (*models.User, error){
	return nil, nil
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error){
	return nil, nil
}