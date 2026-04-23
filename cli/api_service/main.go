package main

import (
	"fmt"
	"log"
	"net/http"

	"go_rest_crud/internal/handler"
	"go_rest_crud/internal/repo/sqlite"
)

/*
TODO LIST:
  Нужно выдавать пользователям выходной документ с табличкой, где будет показано время, кол-во чего-то и т.д
  Нужно это продумать. Но этим я займусь завтра.
  Нужно не тянуть, а то времени немного.
*/

func main() {
	db, dbPath, err := sqlite.OpenDatabase()
	if err != nil {
		log.Fatal("Can't connect to a data base: ", err)
	}
	defer db.Close()
	fmt.Printf("Connected to the SQLite database successfully: %s\n", dbPath)

	var sqliteVersion string

	err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		log.Fatal("Can't print the sqlite version: ", err)
	}
	fmt.Println("SQLite version: ", sqliteVersion)

	if err := sqlite.Migrate(db); err != nil {
		log.Fatal("Can't migrate database schema: ", err)
	}

	sqliteEquipmentStore := sqlite.NewSQLiteEquipmentStore(db)
	sqliteExperimentStore := sqlite.NewSQLiteExperimentStore(db)
	sqliteExperimentEquipmentStore := sqlite.NewSQLiteExperimentEquipmentStore(db)

	equipmentHandler := handler.NewEquipmentHandler(sqliteEquipmentStore)
	experimentHandler := handler.NewExperimentHandler(
		sqliteExperimentStore,
		sqliteEquipmentStore,
		sqliteExperimentEquipmentStore,
	)

	mux := http.NewServeMux()

	mux.Handle("/", &handler.HomeHandler{})
	mux.Handle("/equipment", equipmentHandler)
	mux.Handle("/equipment/", equipmentHandler)

	mux.Handle("/experiment", experimentHandler)
	mux.Handle("/experiment/", experimentHandler)
	mux.Handle("/experiments", experimentHandler)
	mux.Handle("/experiments/", experimentHandler)

	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal("Can't start a server:", err)
	}
}
