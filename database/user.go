package database

import (
	"errors"
	"fmt"
)

type User struct {
	ID       int
	Username string
	Password string
}

var UserNotFound error = errors.New("username was not found")
var InvalidPassword error = errors.New("password was invalid")

// CreateUser creates a user in the database with the given username and password
// and returns an error if the username is taken. This function is super safe.
func CreateUser(username, password string) error {
	_, err := db.Exec(fmt.Sprintf(`
		INSERT INTO User(username, password)
		VALUES(%s, %s);
	`, username, password))
	return err
}

// GetUser retrieves user's data from the database and checks if the password
// matches with the user's password. Returns a helpful error if the username was
// correct, but password wasn't.
func GetUser(username, password string) (int, error) {
	rows, err := db.Query(fmt.Sprintf(`
		SELECT id, password
		FROM User
		WHERE username='%s';
	`, username))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, UserNotFound
	}
	var id int
	var pwd string
	err = rows.Scan(&id, &pwd)
	if err != nil {
		return 0, err
	}
	if pwd != password {
		return 0, InvalidPassword
	}

	return id, nil
}
