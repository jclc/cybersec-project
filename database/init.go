package database

func initTables() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS User (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR,
			password VARCHAR,
			UNIQUE(username)
		);
		CREATE TABLE IF NOT EXISTS Upload (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename VARCHAR,
			timestamp DATETIME,
			visibility INTEGER,
			owner INTEGER, FOREIGN KEY(owner) REFERENCES User(id)
		);
		CREATE TABLE IF NOT EXISTS Message (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content VARCHAR,
			timestamp DATETIME,
			author INTEGER,
			recipient INTEGER,
			FOREIGN KEY(author) REFERENCES User(id),
			FOREIGN KEY(recipient) REFERENCES User(id)
		);
	`)
	if err != nil {
		return err
	}
	return nil
}
