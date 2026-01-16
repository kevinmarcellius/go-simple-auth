package repository

import (
	"gorm.io/gorm"
	model "github.com/kevinmarcellius/go-simple-auth/internal/model"
)

type UserRepository interface {
	// finduserbyid, id is uuid
	FindUserByID(id string) (model.User, error)
	CreateUser(user model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindUserByID implements UserRepository
func (r *userRepository) FindUserByID(id string) (model.User, error) {
	var user model.User
	result := r.db.First(&user, "id = ?", id)
	return user, result.Error
}


func (r *userRepository) CreateUser(user model.User) error {
	result := r.db.Create(&user)
	return result.Error
}
