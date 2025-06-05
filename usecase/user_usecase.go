package usecase

import (
	"product-go-api/model"
	"product-go-api/repository"

	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type UserUsecase struct {
	repository repository.UserRepository
}

func NewUserUsecase(repository repository.UserRepository) UserUsecase {
	return UserUsecase{
		repository: repository,
	}
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
	user, err := uu.repository.GetUserByEmail(req.Email)

	if err != nil || user == nil || user.Password != req.Password {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
