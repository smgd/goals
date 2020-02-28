package store

import (
	"errors"
	"fmt"
	"goals/app/models"

	"golang.org/x/crypto/bcrypt"
)

// User model's repository
type UserRepo struct {
	store *Store
}

func (r *UserRepo) Create(u *models.User) (*models.User, error) {
	if err := r.preCreate(u); err != nil {
		return u, err
	}

	r.store.db.Create(&u)

	return u.Sanitized(), nil
}

func (r *UserRepo) FindByUsername(username string) (*models.User, error) {
	var user models.User

	r.store.db.First(&user, "username = ?", username)

	if user.Username == "" {
		return &user, fmt.Errorf("user with username %s doesn't exists", username)
	}

	return &user, nil
}

func (r *UserRepo) preCreate(u *models.User) error {
	var count int64

	r.store.db.Model(&models.User{}).Where("username = ?", u.Username).Or("email = ?", u.Email).Count(&count)

	if count > 0 {
		return errors.New("user already exists")
	}

	encryptedPassword, err := encryptPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = encryptedPassword

	return nil
}

func encryptPassword(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
