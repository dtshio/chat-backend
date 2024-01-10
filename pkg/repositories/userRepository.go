package repositories

import (
	"fmt"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	core.Repository
}

func (ur *UserRepository) Create(db *gorm.DB, data interface {}) (interface {}, error) {
	user, ok := data.(*models.User)
    if !ok {
        return nil, fmt.Errorf("Data is not a valid user")
    }

	err := db.Create(user).Error
	return *user, err
}

func (ur *UserRepository) FindAll(db *gorm.DB, _ interface {}) (interface {}, error) {
	var users []models.User
	err := db.Find(&users).Error
	return users, err
}

func (ur *UserRepository) FindByID(db *gorm.DB, data interface {}) (interface {}, error) {
	id := data.(uint)
	var user models.User
	err := db.First(&user, id).Error
	return user, err
}

func (ur *UserRepository) FindByEmail(db *gorm.DB, data interface {}) (interface {}, error) {
	email := data.(string)
	var user models.User
	err := db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (ur *UserRepository) Delete(db *gorm.DB, data interface{}) (interface{}, error) {
	id := data.(string)
	err := db.Delete(&models.User{}, id).Error
	return nil, err
}

func (ur *UserRepository) Update(db *gorm.DB, data interface{}) (interface{}, error) {
	user, ok := data.(*models.User)
    if !ok {
        return nil, fmt.Errorf("Data is not a valid user")
    }

	err := db.Save(user).Error
	return *user, err
}

func NewUserRepository() *UserRepository {
	userRepo := &UserRepository{}
	return &UserRepository{
		Repository: *core.NewRepository(
			&models.User{},
			[]core.RepositoryMethod{
				userRepo.Create,
				userRepo.FindAll,
				userRepo.FindByID,
				userRepo.FindByEmail,
				userRepo.Delete,
				userRepo.Update,
			},
		),
	}
}
