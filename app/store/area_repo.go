package store

import (
	"goals/app/models"
)

type AreaRepo struct {
	store *Store
}

func (r *AreaRepo) Create(a *models.Area) (*models.Area, error) {
	r.store.db.Create(&a)

	return a, nil
}

func (r *AreaRepo) FindAreasByUserID(userID uint) (*[]models.Area, error) {
	var areas []models.Area

	r.store.db.Where("user_id = ?", userID).Find(&areas)

	return &areas, nil
}
