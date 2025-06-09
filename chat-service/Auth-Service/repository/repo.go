package repository

import (
	"errors"
	"fmt"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/Auth-Service/kafka"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/Auth-Service/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUser(name, email, password string) (uint, error)
	FindUserByEmail(email string) (*models.User, error)
	FindUserByID(id uint) (*models.User, error)
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
	kp *kafka.KafkaProducer
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB, kafka *kafka.KafkaProducer) UserRepository {
	return &userRepository{db: db,
		kp: kafka}
}

// CreateUser creates a new user with a hashed password
func (r *userRepository) CreateUser(name, email, password string) (uint, error) {
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Create user
	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hash),
	}

	r.kp.KafkaProd(user)
	fmt.Printf("Creating user: %+v\n", user) // Debug: Print user struct
	if err := r.db.Create(&user).Error; err != nil {
		// Check for unique constraint violation (email already exists)
		if errors.Is(err, gorm.ErrDuplicatedKey) || err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return 0, errors.New("email already exists")
		}
		return 0, err
	}

	return user.ID, nil
}

// FindUserByEmail retrieves a user by email
func (r *userRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}

// FindUserByID retrieves a user by ID
func (r *userRepository) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}
