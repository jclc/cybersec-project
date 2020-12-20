package main

import (
	"flag"
	"log"

	"github.com/jclc/cybersec-project/database"
	"github.com/jclc/cybersec-project/server"
)

func main() {
	databaseFile := flag.String("db", "./cybersec.db", "path to the sqlite3 database file")
	port := flag.Int("port", 11037, "HTTP port to use")
	flag.Parse()

	if err := database.Init(*databaseFile); err != nil {
		log.Println("Error initialising databse:", err)
		return
	}
	defer database.Close()

	server.StartServer(*port)
}
