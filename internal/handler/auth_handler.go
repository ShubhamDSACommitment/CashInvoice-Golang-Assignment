package handler

import (
	"net/http"

	"github.com/CashInvoice-Golang-Assignment/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtSecret   string
	jwtExpiry   int
}

func NewAuthHandler(s *service.AuthService, secret string, expiry int) *AuthHandler {
	return &AuthHandler{
		authService: s,
		jwtSecret:   secret,
		jwtExpiry:   expiry,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type RegisterAdminRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email and password required",
		})
		return
	}

	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	token, err := service.GenerateToken(user, h.jwtSecret, h.jwtExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"role":  user.Role,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "valid email and password (min 6 chars) required",
		})
		return
	}

	err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "user already exists or could not create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
	})
}

func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	var req RegisterAdminRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "valid email and password (min 8 chars) required",
		})
		return
	}

	err := h.authService.RegisterAdmin(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "could not create admin user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "admin user created successfully",
	})
}
