package handler

import (
	"errors"
	"net/http"
	"strconv"

	"mini-mall-gin-gorm/internal/service"
	"mini-mall-gin-gorm/pkg/response"

	"github.com/gin-gonic/gin"
)

// ProductHandler 处理商品接口。
type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求参数格式错误")
		return
	}

	product, err := h.productService.Create(c.Request.Context(), req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, ok := parseID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "商品 ID 不合法")
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求参数格式错误")
		return
	}

	product, err := h.productService.Update(c.Request.Context(), id, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, ok := parseID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "商品 ID 不合法")
		return
	}

	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "商品 ID 不合法")
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, product)
}

func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	products, total, err := h.productService.List(c.Request.Context(), page, size)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  products,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

func (h *ProductHandler) handleServiceError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrProductNotFound) {
		response.Fail(c, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, service.ErrInvalidPrice) || errors.Is(err, service.ErrInvalidStock) {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Fail(c, http.StatusInternalServerError, "商品操作失败")
}

func parseID(raw string) (uint, bool) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(id), true
}
