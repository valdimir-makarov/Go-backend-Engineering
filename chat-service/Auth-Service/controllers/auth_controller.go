package controllers

import (
	"log"
	"net/http"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/Auth-Service/kafka"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/Auth-Service/repository"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/Auth-Service/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type Controller struct {
	userRepo  repository.UserRepository
	kafkaProd *kafka.KafkaProducer
}

// NewController creates a new Controller with the given repository
func NewController(userRepo repository.UserRepository, kafka kafka.KafkaProducer) *Controller {
	return &Controller{userRepo: userRepo, kafkaProd: &kafka}
}

func (ctrl *Controller) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Find user
	user, err := ctrl.userRepo.FindUserByEmail(req.Email)
	log.Printf("User found: %+v", user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token,
		"user_id": user})
}

func (ctrl *Controller) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Create user
	_, err := ctrl.userRepo.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "email already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	ctrl.kafkaProd.KafkaProd(req)
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (ctrl *Controller) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID not found"})
		return
	}

	user, err := ctrl.userRepo.FindUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    user.ID,
		"uuid":       user.UUID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

func (ctrl *Controller) GetAllUsers(c *gin.Context) {
	users, err := ctrl.userRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"users": []interface{}{}})
		log.Println("No users found in the database")
		return
	}

	// Log users for debugging
	log.Printf("Retrieved %d users", len(users))
	for _, user := range users {
		log.Printf("User: %+v", user)
	}

	// Send a single JSON response
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
	log.Println("All users retrieved successfully")
}

func (ctrl *Controller) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully. Please discard your JWT.",
	})
}
