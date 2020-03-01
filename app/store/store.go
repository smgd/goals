package store

import (
	"goals/app/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Database wrapper
type Store struct {
	config   *Config
	db       *gorm.DB
	userRepo *UserRepo
	areaRepo *AreaRepo
	goalRepo *GoalRepo
	taskRepo *TaskRepo
}

// Store constructor
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

// Database connection opener wrapper
func (s *Store) Open() error {
	db, err := gorm.Open("postgres", s.config.DatabaseURL)

	if err != nil {
		return err
	}

	db.AutoMigrate(
		&models.User{},
		&models.Area{},
		&models.Goal{},
		&models.Task{},
	)

	s.db = db

	return nil
}

// Database connection closer wrapper
func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) User() *UserRepo {
	if s.userRepo != nil {
		return s.userRepo
	}

	s.userRepo = &UserRepo{
		store: s,
	}

	return s.userRepo
}

func (s *Store) Area() *AreaRepo {
	if s.areaRepo != nil {
		return s.areaRepo
	}

	s.areaRepo = &AreaRepo{
		store: s,
	}

	return s.areaRepo
}

func (s *Store) Goal() *GoalRepo {
	if s.goalRepo != nil {
		return s.goalRepo
	}

	s.goalRepo = &GoalRepo{
		store: s,
	}

	return s.goalRepo
}

func (s *Store) Task() *TaskRepo {
	if s.taskRepo != nil {
		return s.taskRepo
	}

	s.taskRepo = &TaskRepo{
		store: s,
	}

	return s.taskRepo
}
