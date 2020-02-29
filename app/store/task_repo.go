package store

import (
	"goals/app/models"
)

type TaskRepo struct {
	store *Store
}

func (r *TaskRepo) Create(a *models.Task) (*models.Task, error) {
	r.store.db.Create(&a)

	return a, nil
}

func (r *TaskRepo) FindTasksByGoalID(goalID uint) (*[]models.Task, error) {
	var tasks []models.Task

	r.store.db.Where("goal_id = ?", goalID).Find(&tasks)

	return &tasks, nil
}
