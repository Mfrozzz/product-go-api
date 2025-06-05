package repository

import (
	"database/sql"
	"product-go-api/model"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	connection *sql.DB
}

func NewUserRepository(connection *sql.DB) UserRepository {
	return UserRepository{
		connection: connection,
	}
}

func (ur *UserRepository) CreateUser(user model.User) (int, error) {

	var id int

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	query, err := ur.connection.Prepare(
		"INSERT INTO users" + "(username, email, password)" + "VALUES ($1, $2, $3) RETURNING id;",
	)
	if err != nil {
		return 0, err
	}

	err = query.QueryRow(user.Username, user.Email, string(hashedPassword)).Scan(&id)
	if err != nil {
		return 0, err
	}

	query.Close()
	return id, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	query, err := ur.connection.Prepare("SELECT * FROM users WHERE email = $1;")

	if err != nil {
		return nil, err
	}

	var user model.User

	err = query.QueryRow(email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	query.Close()
	return &user, nil
}
