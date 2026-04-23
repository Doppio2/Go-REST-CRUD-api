package main

import (
	"fmt"
	"log"
	"os"

	"go_rest_crud/internal/config"
	"go_rest_crud/internal/repo/sqlite"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Can't load config: ", err)
	}

	command := "up"
	args := os.Args[1:]
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	}

	db, dbPath, err := sqlite.OpenDatabase(cfg.SQLitePath)
	if err != nil {
		log.Fatal("Can't connect to a data base: ", err)
	}
	defer db.Close()

	if err := sqlite.RunMigrationCommand(db, command, args...); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Migration command %q completed for %s\n", command, dbPath)
}
