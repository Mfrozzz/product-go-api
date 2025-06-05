package controller

import (
	"net/http"
	"product-go-api/model"
	"product-go-api/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase usecase.UserUsecase
}

func NewUserController(usecase usecase.UserUsecase) UserController {
	return UserController{
		userUseCase: usecase,
	}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user model.User
	err := ctx.BindJSON(&user)

	if err != nil || user.Username == "" || user.Password == "" {
		response := model.Response{
			Message: "Invalid request body",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	createdUser, err := uc.userUseCase.CreateUser(user)
	if err != nil {
		response := model.Response{
			Message: "Failed to create user",
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

func (uc *UserController) GetUserByEmail(ctx *gin.Context) {
	var req model.LoginRequest
	err := ctx.BindJSON(&req)

	if err != nil || req.Password == "" || req.Email == "" {
		response := model.Response{
			Message: "Invalid request body",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := uc.userUseCase.GetUserByEmail(req)

	if err != nil {
		response := model.Response{
			Message: "Invalid email or password",
		}
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	response := model.Response{
		Message: "Login successful",
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Message": response.Message,
		"token":   token,
	})
}
