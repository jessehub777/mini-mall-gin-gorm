package handler

import (
	"errors"
	"net/http"

	"mini-mall-gin-gorm/internal/middleware"
	"mini-mall-gin-gorm/internal/service"
	"mini-mall-gin-gorm/pkg/response"

	"github.com/gin-gonic/gin"
)

// PurchaseHandler 处理购买相关接口。
type PurchaseHandler struct {
	purchaseService *service.PurchaseService
}

func NewPurchaseHandler(purchaseService *service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{purchaseService: purchaseService}
}

func (h *PurchaseHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	var req service.CreatePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求参数格式错误")
		return
	}

	purchase, err := h.purchaseService.Create(c.Request.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) || errors.Is(err, service.ErrInvalidQuantity) || errors.Is(err, service.ErrInsufficientStock) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "创建购买记录失败")
		return
	}

	response.Success(c, purchase)
}

func (h *PurchaseHandler) ListMine(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	list, err := h.purchaseService.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询购买记录失败")
		return
	}

	response.Success(c, gin.H{"list": list})
}
