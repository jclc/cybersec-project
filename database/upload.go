package database

import (
	"fmt"
	"log"
	"time"
)

type VisibilityState int

const (
	VHidden VisibilityState = iota
	VPublic
)

type Upload struct {
	ID         int64
	Filename   string
	Timestamp  time.Time
	Visibility VisibilityState
	Owner      int64
}

func (u *Upload) Save() error {
	if u.ID == 0 {
		// If ID is 0, create new upload
		timestamp := time.Now().UTC()
		res, err := db.Exec(
			`INSERT INTO Upload(filename, timestamp, visibility, owner) `+
				`VALUES(?, ?, ?, ?);`, u.Filename, timestamp, u.Visibility, u.Owner,
		)
		if err != nil {
			return err
		}
		i, err := res.LastInsertId()
		if err != nil {
			return err
		}
		u.ID = i
	} else {
		// If ID is set, update upload
		_, err := db.Exec(
			`UPDATE Upload `+
				`SET filename=?, visibility=?, owner=?`+
				`WHERE id=?;`, u.Filename, u.Visibility, u.Owner, u.ID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Upload) Delete() error {
	if u.ID == 0 {
		return nil
	}

	_, err := db.Exec(
		`DELETE FROM Upload `+
			`WHERE id=?;`, u.ID,
	)
	if err != nil {
		return err
	}
	u.ID = 0
	return nil
}

func GetUploads(owner int64, onlyPublic bool) []Upload {
	var visReq string
	if onlyPublic {
		visReq = fmt.Sprintf(" AND visibility=%d", VPublic)
	}
	rows, err := db.Query(fmt.Sprintf(
		`SELECT id, filename, timestamp, visibility `+
			`FROM Upload `+
			`WHERE owner=%d%s;`, owner, visReq,
	))
	if err != nil {
		log.Println("Error getting uploads:", err)
		return nil
	}
	defer rows.Close()

	var uploads []Upload
	for rows.Next() {
		var u Upload
		u.Owner = owner
		err := rows.Scan(&u.ID, &u.Filename, &u.Timestamp, &u.Visibility)
		if err != nil {
			log.Println("Error scanning upload", err)
			return nil
		}
		uploads = append(uploads, u)
	}
	return uploads
}

func UploadByID(id int64) (Upload, error) {
	row := db.QueryRow(fmt.Sprintf(
		`SELECT id, filename, timestamp, visibility, owner `+
			`FROM Upload `+
			`WHERE id=%d`, id,
	))

	var u Upload
	err := row.Scan(&u.ID, &u.Filename, &u.Timestamp, &u.Visibility, &u.Owner)
	return u, err
}

func (u *Upload) GetOwner() User {
	user, _ := UserByID(u.Owner)
	return user
}
