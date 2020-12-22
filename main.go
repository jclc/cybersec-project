package main

import (
	"flag"
	"log"

	"github.com/jclc/cybersec-project/database"
	"github.com/jclc/cybersec-project/server"
	"github.com/jclc/cybersec-project/storage"
)

func main() {
	databaseFile := flag.String("db", "./cybersec.db", "path to the sqlite3 database file")
	storagePath := flag.String("storage", "", "directory to use for file storage (default is in the system temp folder)")
	sessionKey := flag.String("session-key", "ASD87B90F82CA035E7FA", "session key for cookie storage")
	port := flag.Int("port", 11037, "HTTP port to use")
	flag.Parse()

	if err := storage.Init(*storagePath); err != nil {
		log.Println("Error initalising storage:", err)
		return
	}
	defer storage.Close()

	if err := database.Init(*databaseFile); err != nil {
		log.Println("Error initialising databse:", err)
		return
	}
	defer database.Close()
	database.CreateTestData()

	server.StartServer(*port, *sessionKey)
}
