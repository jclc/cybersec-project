package database

import (
	"fmt"
	"time"
)

type Message struct {
	ID        int64
	Content   string
	Timestamp time.Time
	Author    string
	Recipient int64
}

func (m *Message) Save() error {
	if m.ID == 0 {
		// If ID is 0, create new message
		timestamp := time.Now().UTC()
		res, err := db.Exec(
			`INSERT INTO Message(content, timestamp, author, recipient) `+
				`VALUES(?, ?, ?, ?);`, m.Content, timestamp, m.Author, m.Recipient,
		)
		if err != nil {
			return err
		}
		i, err := res.LastInsertId()
		if err != nil {
			return err
		}
		m.ID = i
	} else {
		// If ID is set, update user
		_, err := db.Exec(
			`UPDATE Message `+
				`SET content=?, timestamp=?, author=?, recipient=?`+
				`WHERE id=?;`, m.Content, m.Timestamp, m.Author, m.Recipient, m.ID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func MessageByID(id int64) (Message, error) {
	row := db.QueryRow(fmt.Sprintf(
		`SELECT id, content, timestamp, author, recipient `+
			`FROM Message `+
			`WHERE id=%d;`, id,
	))

	var m Message
	err := row.Scan(&m.ID, &m.Content, &m.Timestamp, &m.Author, &m.Recipient)
	return m, err
}
