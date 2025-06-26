package controller

import (
	"net/http"
	"product-go-api/model"
	"product-go-api/usecase"
	"strconv"

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

	if user.Role == "" {
		user.Role = "user"
	}

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

	response := model.Response{
		Message: "Register successful",
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"Message": response.Message,
		"User": map[string]interface{}{
			"user_id":  createdUser.ID,
			"email":    createdUser.Email,
			"username": createdUser.Username,
			"role":     createdUser.Role,
		},
	})
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

func (uc *UserController) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id_user")

	if id == "" {
		response := model.Response{
			Message: "id_user is required",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	id_user, err := strconv.Atoi(id)
	if err != nil || id_user < 1 {
		response := model.Response{
			Message: "id_user must be a positive number",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := uc.userUseCase.GetUserById(id_user)
	if err != nil {
		response := model.Response{
			Message: "Failed to retrieve user",
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if user == nil {
		response := model.Response{
			Message: "User not found",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetUserInfo(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDValue.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in context"})
		return
	}

	user, err := uc.userUseCase.GetUserById(userID)
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id_user")

	if id == "" {
		response := model.Response{
			Message: "id_user is required",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	id_user, err := strconv.Atoi(id)
	if err != nil || id_user < 1 {
		response := model.Response{
			Message: "id_user must be a positive number",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = uc.userUseCase.DeleteUser(id_user)
	if err != nil {
		response := model.Response{
			Message: "Failed to delete user.",
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := model.Response{
		Message: "User deleted successfully",
	}
	ctx.JSON(http.StatusOK, response)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id_user")

	if id == "" {
		response := model.Response{
			Message: "id_user is required",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	id_user, err := strconv.Atoi(id)
	if err != nil || id_user < 1 {
		response := model.Response{
			Message: "id_user must be a positive number",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	existingUser, err := uc.userUseCase.GetUserById(id_user)
	if err != nil {
		response := model.Response{
			Message: "Failed to retrieve user.",
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if existingUser == nil {
		response := model.Response{
			Message: "User not found",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	var updateData map[string]interface{}
	if err := ctx.BindJSON(&updateData); err != nil {
		response := model.Response{
			Message: "Invalid JSON.",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if username, ok := updateData["username"].(string); ok {
		existingUser.Username = username
	}
	if email, ok := updateData["email"].(string); ok {
		existingUser.Email = email
	}
	if password, ok := updateData["password"].(string); ok && password != "" {
		existingUser.Password = password
	}
	if newRole, ok := updateData["role"].(string); ok {
		if role, exists := ctx.Get("role"); !exists || role != "admin" {
			response := model.Response{
				Message: "Only admin can change user role.",
			}
			ctx.JSON(http.StatusForbidden, response)
			return
		}
		existingUser.Role = newRole
	}

	updatedUser, err := uc.userUseCase.UpdateUser(*existingUser)
	if err != nil {
		response := model.Response{
			Message: "Failed to update user.",
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}
