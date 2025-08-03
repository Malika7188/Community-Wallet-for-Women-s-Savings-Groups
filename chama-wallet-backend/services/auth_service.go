package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/utils"
)

var jwtSecret = []byte("your-secret-key-change-in-production") // Change this in production

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(userID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RegisterUser creates a new user account
func RegisterUser(req models.RegisterRequest) (models.AuthResponse, error) {
	// Generate wallet
	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return models.AuthResponse{}, errors.New("user with this email already exists")
	}
	
	wallet, err := utils.GenerateStellarWallet()
	if err != nil {
		return models.AuthResponse{}, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return models.AuthResponse{}, err
	}

	// Create user with secret key
	user := models.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Wallet:    wallet.PublicKey,
		SecretKey: wallet.SecretKey, // Make sure this is set
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return models.AuthResponse{}, err
	}

	// Generate token
	token, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// LoginUser authenticates a user and returns a token
func LoginUser(req models.LoginRequest) (models.AuthResponse, error) {
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return models.AuthResponse{}, errors.New("invalid email or password")
	}

	// Check password
	if !CheckPasswordHash(req.Password, user.Password) {
		return models.AuthResponse{}, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID string) (models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

// UpdateUser updates a user's information
func UpdateUser(user models.User) error {
	return database.DB.Save(&user).Error
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}
