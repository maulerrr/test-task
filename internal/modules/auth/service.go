package auth

import (
	"fmt"
	"test-task/internal/config"
	db "test-task/internal/database"
	"test-task/internal/modules/auth/models"
	"test-task/pkg/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Handler db.DBHandler
	Config  *config.Config
}

func InitAuthService(handler db.DBHandler, cfg *config.Config) (*Service, error) {
	return &Service{
		Handler: handler,
		Config:  cfg,
	}, nil
}

func (s *Service) CreateUser(email, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.Handler.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AuthenticateUser(email, password string) (uuid.UUID, error) {
	var user models.User
	if err := s.Handler.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return uuid.Nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (s *Service) SaveRefreshToken(userID uuid.UUID, refreshToken, ipAddress string) error {
	if err := s.Handler.DB.Where("user_id = ?", userID).Delete(&models.Token{}).Error; err != nil {
		return err
	}

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	token := &models.Token{
		ID:               uuid.New(),
		UserID:           userID,
		RefreshTokenHash: string(hashedToken),
		IPAddress:        ipAddress,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(30 * 24 * time.Hour),
	}

	return s.Handler.DB.Create(token).Error
}

func (s *Service) ValidateRefreshToken(refreshToken, accessToken, ipAddress string) (bool, error) {
	userIDStr, err := utils.ExtractUserIDFromToken(accessToken, s.Config.JWTSecretKey)
	if err != nil {
		return false, err
	}

	userID, err := utils.ConvertStringToUUID(userIDStr)
	if err != nil {
		return false, err
	}

	var token models.Token
	if err := s.Handler.DB.Where("user_id = ? AND expires_at > NOW()", userID).First(&token).Error; err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(token.RefreshTokenHash), []byte(refreshToken)); err != nil {
		return false, err
	}

	user, err := s.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	if token.IPAddress != ipAddress {
		s.MockSendEmailWarning(user.Email, token.IPAddress, ipAddress)
		return false, fmt.Errorf("IP address mismatch")
	}

	return true, nil
}

func (s *Service) CleanupExpiredTokens() error {
	return s.Handler.DB.Where("expires_at < NOW()").Delete(&models.Token{}).Error
}

func (s *Service) MockSendEmailWarning(email, oldIP, newIP string) {
	fmt.Printf("WARNING: IP address change detected!\n")
	fmt.Printf("Email sent to: %s\n", email)
	fmt.Printf("Old IP: %s\n", oldIP)
	fmt.Printf("New IP: %s\n", newIP)
}

func (s *Service) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.Handler.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
