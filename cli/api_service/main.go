package main

import (
	"log"
	"fmt"
	"net/http"
	"database/sql"

    _ "github.com/glebarez/go-sqlite"
	
	"go_rest_crud/internal/handler"
	"go_rest_crud/internal/repo"
)

func main() {
// -- connecting sqlite -- 
	/*sql.Open("sqlite", "db/equipment.db")*/
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal("Can't connect to a data base: ", err)
	}
	defer db.Close()
    fmt.Println("Connected to the SQLite database successfully.")

	var sqliteVersion string

	err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		log.Fatal("Can't print the sqlite version: ", err)
	}
	fmt.Println("SQLite version: ", sqliteVersion)

	sqlCreateTable := `
	CREATE TABLE IF NOT EXISTS equipment (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		field TEXT
	);`
	_, err = db.Exec(sqlCreateTable)
	if err != nil {
		log.Fatal("Can't create a table: ", err)
	}
    fmt.Println("Table \"equipment\" was created successfully.")

	sqlite_store := repo.NewSQLiteStore(db)

	// NOTE: test
	fmt.Println("Schema of table equipment: ")
	sqlite_store.PrintDBSchema()
	// -- -- 
	equipmentHandler := handler.NewEquipmentHandler(sqlite_store)

	mux := http.NewServeMux()
	
	mux.Handle("/", &handler.HomeHandler{})
	mux.Handle("/equipment", equipmentHandler)
	mux.Handle("/equipment/", equipmentHandler)

	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal("Can't start a server:", err)
	}
}
