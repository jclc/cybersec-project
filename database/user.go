package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type User struct {
	ID             int64
	Username       string
	Password       string
	SocialSecurity string
	Uploads        int // Number of public uploads when querying
}

var UserNotFound error = errors.New("username was not found")
var InvalidPassword error = errors.New("password was invalid")

// CreateUser creates a user in the database with the given username and password
// and returns an error if the username is taken. This function is super safe.
func (u *User) Create() error {
	res, err := db.Exec(fmt.Sprintf(
		`INSERT INTO User(username, password, social_security) `+
			`VALUES('%s', '%s', '%s');`, u.Username, u.Password, u.SocialSecurity))
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id
	return nil
}

func (u *User) GetMessages() []Message {
	rows, err := db.Query(
		`SELECT id, content, timestamp, author, recipient `+
			`FROM Messages `+
			`WHERE recipient=?`, u.ID,
	)
	if err != nil {
		log.Println("Error getting messages:", err)
		return nil
	}
	var messages []Message
	for rows.Next() {
		var msg Message
		err = rows.Scan(&msg.ID, &msg.Content, &msg.Timestamp, &msg.Author, &msg.Recipient)
		if err != nil {
			log.Println("Error scanning messages:", err)
			return nil
		}
		messages = append(messages, msg)
	}
	return messages
}

// GetUser retrieves user's data from the database and checks if the password
// matches with the user's password. Returns a helpful error if the username was
// correct, but password wasn't. This function is super safe.
func GetUser(username, password string) (User, error) {
	row := db.QueryRow(fmt.Sprintf(
		`SELECT id, username, password, social_security `+
			`FROM User `+
			`WHERE username='%s';`, username,
	))

	var u User
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.SocialSecurity,
	); err == sql.ErrNoRows {
		return u, UserNotFound
	} else if err != nil {
		return u, err
	}
	if u.Password != password {
		return User{}, InvalidPassword
	}

	return u, nil
}

func GetUsers() []User {
	rows, err := db.Query(fmt.Sprintf(
		`SELECT id, username ` +
			`from User;`,
	))
	if err != nil {
		log.Println("Error getting users:", err)
		return nil
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username)
		if err != nil {
			log.Println("Error scanning user", err)
			return nil
		}
		users = append(users, u)
	}
	return users
}

func UserByID(id int64) (User, error) {
	row := db.QueryRow(fmt.Sprintf(
		`SELECT id, username, password, social_security `+
			`FROM User `+
			`WHERE id=%d`, id,
	))

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.SocialSecurity)
	return u, err
}
