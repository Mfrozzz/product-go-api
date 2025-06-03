package controller

import (
	"net/http"
	"product-go-api/model"
	"product-go-api/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productController struct {
	productUseCase usecase.ProductUsecase
}

func NewProductController(usecase usecase.ProductUsecase) productController {
	return productController{
		productUseCase: usecase,
	}
}

func (p *productController) GetProducts(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "page must be a positive number"})
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a positive number"})
		return
	}
	name := ctx.Query("name")

	products, err := p.productUseCase.GetProducts(page, limit, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, products)
}

func (p *productController) CreateProduct(ctx *gin.Context) {
	var product model.Product
	err := ctx.BindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if product.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product name is required"})
		return
	}

	if product.Price < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Price must be non-negative"})
		return
	}

	insertedProduct, err := p.productUseCase.CreateProduct(product)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, insertedProduct)
}

func (p *productController) GetProductById(ctx *gin.Context) {

	id := ctx.Param("id_product")

	if id == "" {
		response := model.Response{
			Message: "id_product is required",
		}

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	id_product, err := strconv.Atoi(id)

	if id_product < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id_product must be a positive number"})
		return
	}

	if err != nil {
		response := model.Response{
			Message: "id_product must be a number",
		}

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := p.productUseCase.GetProductById(id_product)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	if product == nil {
		response := model.Response{
			Message: "Product not found",
		}

		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (p *productController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id_product")

	if id == "" {
		response := model.Response{
			Message: "id_product is required",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	id_product, err := strconv.Atoi(id)

	if id_product < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id_product must be a positive number"})
		return
	}

	if err != nil {
		response := model.Response{
			Message: "id_product must be a number",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = p.productUseCase.DeleteProduct(id_product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response := model.Response{
		Message: "Product deleted successfully",
	}
	ctx.JSON(http.StatusOK, response)
}

func (p *productController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id_product")

	if id == "" {
		response := model.Response{
			Message: "id_product is required",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	id_product, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{
			Message: "id_product must be a number",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	existingProduct, err := p.productUseCase.GetProductById(id_product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingProduct == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var updateData map[string]interface{}
	if err := ctx.BindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if name, ok := updateData["name"].(string); ok {
		existingProduct.Name = name
	}

	if price, ok := updateData["price"].(float64); ok {
		existingProduct.Price = price
	}

	updatedProduct, err := p.productUseCase.UpdateProduct(*existingProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)

}
