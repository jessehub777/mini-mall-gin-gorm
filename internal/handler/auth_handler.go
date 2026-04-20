package handler

import (
	"errors"
	"net/http"

	"mini-mall-gin-gorm/internal/service"
	"mini-mall-gin-gorm/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler 处理注册和登录接口。
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求参数格式错误")
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "注册失败")
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求参数格式错误")
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Fail(c, http.StatusUnauthorized, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "登录失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
		},
	})
}
