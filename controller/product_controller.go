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
	products, err := p.productUseCase.GetProducts()
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

	insertedProduct, err := p.productUseCase.CreateProduct(product)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
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

	var product model.Product
	err = ctx.BindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	product.ID = id_product

	updatedProduct, err := p.productUseCase.UpdateProduct(product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)

}
