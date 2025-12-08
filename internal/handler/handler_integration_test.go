package handler

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io"
	"os"
	"bytes"
	"database/sql"
	"log"
	"fmt"

    _ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/gosimple/slug"

	"go_rest_crud/internal/repo"
)

// Функция для тестов. Читает тестовый набор данных.
func readTestData(t *testing.T, name string) []byte {
    t.Helper()
    content, err := os.ReadFile("testdata/" + name)
    if err != nil {
        t.Errorf("Could not read %v", name)
    }

    return content
}

// Функция для тестов. Инициализирует sqlite.
func InitSqliteTest() *repo.SQLiteStore {
	// TODO: Позже нужно будет создавать файл с бд.
	// Пока что я использую :memory параметр, что хранит всю бд в оперативной памяти.
	// sql.Open("sqlite", "db/equipment.db")
	// sql.Open("sqlite", "db/equipment.db")
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal("Can't connect to a data base: ", err)
	}
    fmt.Println("Connected to the SQLite database successfully.")

	s := repo.NewSQLiteStore(db)

	return s
}

// TODO: Нужно вынести эту функцию в отдельный пакет какой-то.
// Функция для тестов. Создаем таблицу equipment для тестирования работоспособности.
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

// Функция для тестов. Основная тестирующая функция, которая проверяет работоспособность всего веб-приложения.
// TODO: Переписать тесты.
func TestEquipmentHandlerCRUD_Integration(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal("Can't connect to a data base: ", err)
	}
	defer db.Close()
    fmt.Println("Connected to the SQLite database successfully.")

	CreateTables(db);

	sqliteEquipmentStore := sqlite.NewSQLiteEquipmentStore(db)
	equipmentHandler = NewEquipmentHandler(sqliteEquipmentStore);

	testData1 := readTestData(t, "spectrophotometer.json")
	testData2 := readTestData(t, "high_speed_centrifuge.json")
	testData1Reader := bytes.NewReader(testData1)
	testData2Reader := bytes.NewReader(testData2)

	// -- CREATE: POST /equipment --
	r := httptest.NewRequest(http.MethodPost, "/equipment", testData1Reader)
	w := httptest.NewRecorder()
	
	/*
	sqliteStore := InitSqliteTest()
	SQLiteCreateEquipmentTableTest(sqliteStore)
	defer sqliteStore.Close()
	
	// --- Test table schema ---
	fmt.Println("Schema of table equipment: ")
	sqliteStore.PrintDBSchema()

	equipmentHandler := NewEquipmentHandler(sqliteStore)

	equipmentTestData1 := readTestData(t, "spectrophotometer.json")
	equipmentTestDataReader1 := bytes.NewReader(equipmentTestData1)

	equipmentTestData2 := readTestData(t, "high_speed_centrifuge.json")
	//equipmentTestDataReader2 := bytes.NewReader(equipmentTestData2)

	// --- CREATE: POST /equipment ---
	r := httptest.NewRequest(http.MethodPost, "/equipment", equipmentTestDataReader1)
	w := httptest.NewRecorder()
	equipmentHandler.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	saved, _ := sqliteStore.List()
	assert.Len(t, saved, 1)

	createdSlug := slug.Make("Spectrophotometer") // "spectrophotometer"

	// --- LIST: GET /equipment/ ---
	r = httptest.NewRequest(http.MethodGet, "/equipment/", nil)
	w = httptest.NewRecorder()
	equipmentHandler.ServeHTTP(w, r)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	// Correction for handler_equipment_test.go:105:
	// The expected JSON must include "id":1 (since ID is not zero) and use the slug for "name".
	expectedList := fmt.Sprintf(`[{"id":1,"name":"%s","field":"Chemistry / Biochemistry"}]`, createdSlug)
	assert.JSONEq(t, expectedList, string(data))

	// --- UPDATE: PUT /equipment/spectrophotometer ---
	// Сброс ридера для повторного чтения equipmentTestData2
	// Важно: в тесте был вызов Reset() для *bytes.Reader, но если ридер был исчерпан,
	// необходимо убедиться, что он сброшен корректно.
	// Лучше создать новый ридер, чтобы избежать проблем с состоянием.
	equipmentTestDataReaderForUpdate := bytes.NewReader(equipmentTestData2)
	
	r = httptest.NewRequest(http.MethodPut, "/equipment/"+createdSlug, equipmentTestDataReaderForUpdate)
	w = httptest.NewRecorder()
	equipmentHandler.ServeHTTP(w, r)

	res = w.Result()
	defer res.Body.Close()
	// Проверка статуса: должен быть 200 OK при успешном обновлении
	assert.Equal(t, http.StatusOK, res.StatusCode) 

	// New slug after update (e.g., "high-speed-centrifuge"), used for subsequent GET/DELETE
	newSlug := slug.Make("High-Speed Centrifuge")

	// Correction for handler_equipment_test.go:123: Get must use the new slug, as the key in the DB was updated
	updatedEquipment, err := sqliteStore.Get(newSlug)
	assert.NoError(t, err)
	// Correction: The stored Name is the slug
	assert.Equal(t, newSlug, updatedEquipment.Name)
	assert.Equal(t, "Biology / Molecular Genetics", updatedEquipment.Field)

	// --- DELETE: DELETE /equipment/spectrophotometer ---
	// Correction for handler_equipment_test.go:135: Delete must use the new slug
	r = httptest.NewRequest(http.MethodDelete, "/equipment/"+newSlug, nil)
	w = httptest.NewRecorder()
	equipmentHandler.ServeHTTP(w, r)

	res = w.Result()
	defer res.Body.Close()
	// Ожидаемый статус: 204 No Content
	assert.Equal(t, http.StatusNoContent, res.StatusCode) // ← Ожидаем 204

	savedAfterDelete, _ := sqliteStore.List()
	// Ожидаем 0 элементов после удаления
	assert.Len(t, savedAfterDelete, 0)
	*/
}

func TestExperimentHandlerCRUD_Integration(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal("Can't connect to a data base: ", err)
	}
	defer db.Close()
    fmt.Println("Connected to the SQLite database successfully.")

	CreateTables(db);

	sqliteEquipmentStore := sqlite.NewSQLiteEquipmentStore(db)
	sqliteExperimentStore := sqlite.NewSQLiteExperimentStore(db)
	sqliteExperimentEquipmentStore := sqlite.NewSQLiteExperimentEquipmentStore(db)
}
