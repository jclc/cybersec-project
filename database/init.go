package database

import (
	"fmt"
	"log"
)

func initTables() error {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS User ( ` +
			`id INTEGER PRIMARY KEY AUTOINCREMENT, ` +
			`username VARCHAR, ` +
			`password VARCHAR, ` +
			`social_security VARCHAR, ` +
			`UNIQUE(username) ` +
			`); ` +
			`CREATE TABLE IF NOT EXISTS Upload ( ` +
			`id INTEGER PRIMARY KEY AUTOINCREMENT, ` +
			`filename VARCHAR, ` +
			`timestamp DATETIME DEFAULT CURRENT_TIMESTAMP, ` +
			`visibility INTEGER, ` +
			`owner INTEGER, FOREIGN KEY(owner) REFERENCES User(id) ` +
			`); ` +
			`CREATE TABLE IF NOT EXISTS Message ( ` +
			`id INTEGER PRIMARY KEY AUTOINCREMENT, ` +
			`content VARCHAR, ` +
			`timestamp DATETIME DEFAULT CURRENT_TIMESTAMP, ` +
			`author VARCHAR, ` +
			`recipient INTEGER, ` +
			`FOREIGN KEY(recipient) REFERENCES User(id) ` +
			`);`,
	)
	if err != nil {
		return err
	}
	return nil
}

func CreateTestData() {
	log.Println("Creating test data")
	usr := User{0, "admin", "admin", "777", 0}
	err := usr.Create()
	if err != nil {
		log.Println(err)
		log.Println("Admin already exists")
		return
	}

	admin, err := GetUser("admin", "admin")
	if err != nil {
		log.Println("Error getting admin:", err)
		return
	}

	for i := 0; i < 10; i++ {
		var upload Upload
		upload.Filename = fmt.Sprintf("file%d.bin", i)
		upload.Owner = admin.ID
		upload.Visibility = VisibilityState(i % 2)
		err := upload.Save()
		if err != nil {
			log.Println("Error saving upload:", err)
			return
		}
	}

	for i := 0; i < 10; i++ {
		var message Message
		message.Content = []byte(fmt.Sprintf("This is the message #%d", i))
		message.Author = admin.ID
		message.Recipient = admin.ID
		err := message.Save()
		if err != nil {
			log.Println("Error saving message:", err)
			return
		}
	}
}
