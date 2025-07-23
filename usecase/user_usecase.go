package usecase

import (
	"product-go-api/model"
	"product-go-api/repository"

	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repository repository.UserRepository
}

func NewUserUsecase(repository repository.UserRepository) UserUsecase {
	return UserUsecase{
		repository: repository,
	}
}

func (uu *UserUsecase) GetUsers(page, limit int, name string) ([]model.User, error) {
	return uu.repository.GetUsers(page, limit, name)
}

func (uu *UserUsecase) CreateUser(user model.User) (model.User, error) {
	existingUser, err := uu.repository.GetUserByEmail(user.Email)

	if err != nil {
		return model.User{}, err
	}

	if existingUser != nil {
		return *existingUser, nil
	}

	userId, err := uu.repository.CreateUser(user)
	if err != nil {
		return model.User{}, err
	}

	user.ID = userId
	return user, nil
}

func (uu *UserUsecase) GetUserByEmail(req model.LoginRequest) (string, error) {
	godotenv.Load()
	var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	user, err := uu.repository.GetUserByEmail(req.Email)

	if err != nil || user == nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (uu *UserUsecase) GetUserById(id_user int) (*model.User, error) {
	user, err := uu.repository.GetUserById(id_user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uu *UserUsecase) DeleteUser(id_user int) error {
	err := uu.repository.DeleteUser(id_user)
	if err != nil {
		return err
	}
	return nil
}

func (uu *UserUsecase) UpdateUser(user model.User) (model.User, error) {
	updatedUser, err := uu.repository.UpdateUser(user)
	if err != nil {
		return model.User{}, err
	}
	return *updatedUser, nil
}
