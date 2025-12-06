package main

import (
	"log"
	"fmt"
	"net/http"
	"database/sql"

    _ "github.com/glebarez/go-sqlite"
	
	"go_rest_crud/internal/handler"
//  "go_rest_crud/internal/repo"
	"go_rest_crud/internal/repo/sqlite"
)

// TODO: это нужно перенести в файлы миграции (Я ток не знаю, можно ли так делать в sqlite или это внутри кода надо делать) sql.
// Либо в отдельный файл .go файл.
func CreateTables(db *sql.DB) {
	// Таблица для оборудования.
	sqlCreateTable := `
	CREATE TABLE IF NOT EXISTS equipment (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT
	);`
	_, err := db.Exec(sqlCreateTable)
	if err != nil {
		log.Fatal("Can't create table \"equipment\": ", err)
	}
	fmt.Println("Table \"equipment\" was created successfully.")

	// Таблица для эксперимента.
	sqlCreateTable = `
	CREATE TABLE IF NOT EXISTS experiment (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT
	);`
	_, err = db.Exec(sqlCreateTable)
	if err != nil {
		log.Fatal("Can't create table experiment: ", err)
	}
	fmt.Println("Table \"experiment\" was created successfully.")

	// Объединенная таблица.
	sqlCreateTable = `
	CREATE TABLE IF NOT EXISTS equipment_experiment (
		experiment_id INTEGER NOT NULL,
		equipment_id INTEGER NOT NULL,
		PRIMARY KEY (experiment_id, equipment_id),
		FOREIGN KEY (experiment_id) REFERENCES experiment(id),
		FOREIGN KEY (equipment_id) REFERENCES equipment(id)
	);`
	_, err = db.Exec(sqlCreateTable)
	if err != nil {
		log.Fatal("Can't create table equipment_experiment: ", err)
	}
	fmt.Println("Table \"equipment_experiment\" was created successfully.")
}

func main() {
	// -- Подключение базы данных sqlite. -- 
	// TODO: Позже нужно будет создавать файл с бд.
	// Пока что я использую :memory параметр, что хранит всю бд в оперативной памяти.
	// sql.Open("sqlite", "db/equipment.db")
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal("Can't connect to a data base: ", err)
	}
	defer db.Close()
    fmt.Println("Connected to the SQLite database successfully.")
	// ----

	// -- Проверка работоспособности пакета и вывод версии. --
	// NOTE: наверное, в релизной версии стоит это убрать.
	var sqliteVersion string

	err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		log.Fatal("Can't print the sqlite version: ", err)
	}
	fmt.Println("SQLite version: ", sqliteVersion)

	CreateTables(db) // Создание необходимых таблиц.
	// ----

	// -- Создание хранилищ для операций с БД --.
	// Repo - интерфейс содержащий все операции для работы с бд.
	sqliteEquipmentStore := sqlite.NewSQLiteEquipmentStore(db)
	sqliteExperimentStore := sqlite.NewSQLiteExperimentStore(db)
	// TODO: Я не знаю, нужно ли создавать отдельную таблицу для связи.
	// Скорее всего это плохая идея, но кто его знает. Мне может понадобится, если я захочу еще какое-нибудь поле добавить.
	// Поэтому пускай будет. Это довольно легко можно будет соединить с хранилищем ExperimentStore.
	sqliteExperimentEquipmentStore := sqlite.NewSQLiteExperimentEquipmentStore(db)
	// --- 

	// -- Создаем ручку и подключение к http серверу и связываем все это между собой. --
	// TODO: сделать ручки для остальных таблиц.
	// NOTE: Я думаю, что для ExperimentEquipment мне не нужен отдельный handler,
	// а можно использовать вложенные маршруты. 
	equipmentHandler := handler.NewEquipmentHandler(sqliteEquipmentStore)
	experimentHandler := handler.NewExperimentHandler(sqliteExperimentStore, sqliteExperimentEquipmentStore)

	mux := http.NewServeMux()
	
	mux.Handle("/", &handler.HomeHandler{})
	mux.Handle("/equipment", equipmentHandler)
	mux.Handle("/equipment/", equipmentHandler)

	// TODO: реализовать оставшиеся ручки.
	mux.Handle("/experiments", experimentHandler)            // Для взаимодействия со всем списком экспериментов.
	mux.Handle("/experiment/", experimentHandler)            // Для взаимодействия со определенным экспериментом.

	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal("Can't start a server:", err)
	}
	// ----
}
