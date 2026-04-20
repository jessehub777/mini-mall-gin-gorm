package handler

import (
	"net/http"

	"mini-mall-gin-gorm/internal/middleware"
	"mini-mall-gin-gorm/internal/service"
	"mini-mall-gin-gorm/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserHandler 处理用户个人信息接口。
type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取个人信息失败")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求参数格式错误")
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新个人信息失败")
		return
	}

	response.Success(c, user)
}
