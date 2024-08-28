package auth

import (
	"net/http"
	"test-task/internal/config"
	"test-task/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
	Config  *config.Config
}

func NewHandler(service *Service, cfg *config.Config) *Handler {
	return &Handler{
		Service: service,
		Config:  cfg,
	}
}

func (h *Handler) RegisterUserHandler(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.Service.CreateUser(requestBody.Email, requestBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	ipAddress := c.ClientIP()
	accessToken, err := utils.GenerateAccessToken(user.ID.String(), ipAddress, h.Config.JWTSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate refresh token"})
		return
	}

	if err := h.Service.SaveRefreshToken(user.ID, refreshToken, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) LoginUserHandler(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userID, err := h.Service.AuthenticateUser(requestBody.Email, requestBody.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	ipAddress := c.ClientIP()
	accessToken, err := utils.GenerateAccessToken(userID.String(), ipAddress, h.Config.JWTSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate refresh token"})
		return
	}

	if err := h.Service.SaveRefreshToken(userID, refreshToken, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) IssueTokensHandler(c *gin.Context) {
	userIDStr := c.Param("id")
	ipAddress := c.ClientIP()

	accessToken, err := utils.GenerateAccessToken(userIDStr, ipAddress, h.Config.JWTSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate refresh token"})
		return
	}

	userID, err := utils.ConvertStringToUUID(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error occurred"})
		return
	}
	if err := h.Service.SaveRefreshToken(userID, refreshToken, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) RefreshTokensHandler(c *gin.Context) {
	var requestBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ipAddress := c.ClientIP()

	valid, err := h.Service.ValidateRefreshToken(requestBody.RefreshToken, requestBody.AccessToken, ipAddress)
	if err != nil {
		if err.Error() == "IP address mismatch" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "IP address mismatch"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token or IP address mismatch"})
		return
	}

	userIDStr, err := utils.ExtractUserIDFromToken(requestBody.AccessToken, h.Config.JWTSecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(userIDStr, ipAddress, h.Config.JWTSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate refresh token"})
		return
	}

	userID, err := utils.ConvertStringToUUID(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error occurred"})
		return
	}

	if err := h.Service.SaveRefreshToken(userID, newRefreshToken, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
