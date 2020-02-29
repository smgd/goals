package store

import (
	"goals/app/models"
)

type GoalRepo struct {
	store *Store
}

func (r *GoalRepo) Create(a *models.Goal) (*models.Goal, error) {
	r.store.db.Create(&a)

	return a, nil
}

func (r *GoalRepo) FindGoalsByAreaID(areaID uint) (*[]models.Goal, error) {
	var goals []models.Goal

	r.store.db.Where("area_id = ?", areaID).Find(&goals)

	return &goals, nil
}
