package handler_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/assert"

	"go_rest_crud/internal/entity"
	"go_rest_crud/internal/handler"
	"go_rest_crud/internal/repo/sqlite"
)

// TODO: Переписать тесты нормально.
// Сейчас это нейро-слоп. Хотя мб и заработает.
// Но мне кажется проще тесты самому написать.

// ErrorResponse используется для десериализации ответа об ошибке из handler_error.go
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

// ====================================================================================================
//                                    НАСТРОЙКА ИНТЕГРАЦИОННОЙ СРЕДЫ
// ====================================================================================================

// setupTestDB создает in-memory SQLite DB, инициализирует необходимые таблицы
// и настраивает роутинг в точности как в вашем main-файле.
func setupTestDB(t *testing.T) (*sql.DB, *http.ServeMux) {
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	assert.NoError(t, err)

	err = sqlite.Migrate(db)
	assert.NoError(t, err)

	// Инициализация хранилищ
	equipmentStore := sqlite.NewSQLiteEquipmentStore(db)
	experimentStore := sqlite.NewSQLiteExperimentStore(db)
	experimentEquipmentStore := sqlite.NewSQLiteExperimentEquipmentStore(db)

	// Инициализация хэндлеров
	equipmentHandler := handler.NewEquipmentHandler(equipmentStore)
	experimentHandler := handler.NewExperimentHandler(
		experimentStore,
		equipmentStore,
		experimentEquipmentStore,
	)

	// Настройка роутера, полностью имитирующая ваш main-файл
	mux := http.NewServeMux()

	mux.Handle("/", &handler.HomeHandler{})

	// Оборудование
	mux.Handle("/equipment", equipmentHandler)
	mux.Handle("/equipment/", equipmentHandler) // Для GET/PUT/DELETE /equipment/{id}

	mux.Handle("/experiment", experimentHandler)
	mux.Handle("/experiment/", experimentHandler)
	mux.Handle("/experiments", experimentHandler)
	mux.Handle("/experiments/", experimentHandler)

	return db, mux
}

// readTestData читает JSON-файл из папки testdata.
func readTestData(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + filename)
	assert.NoError(t, err, "Не удалось прочитать файл %s. Убедитесь, что папка 'testdata' существует и содержит файл.", filename)
	return data
}

// executeRequest - вспомогательная функция для выполнения HTTP-запроса
func executeRequest(r *http.ServeMux, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

// ====================================================================================================
//                                     ТЕСТЫ ДЛЯ EQUIPMENT
// ====================================================================================================

func TestEquipmentHandlerCRUD_Integration(t *testing.T) {
	db, mux := setupTestDB(t)
	defer db.Close()

	var createdID int

	// 1. POST /equipment (Create/Add)
	t.Run("Equipment_Create_Success", func(t *testing.T) {
		payload := readTestData(t, "new_equipment.json")
		req, _ := http.NewRequest("POST", "/equipment", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(mux, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response SuccessResponse[entity.Equipment]
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Data.ID > 0, "Ожидался ID")
		createdID = response.Data.ID // Сохраняем ID
	})

	// 2. GET /equipment/{id} (Read/Get)
	t.Run("Equipment_Get_Success", func(t *testing.T) {
		url := fmt.Sprintf("/equipment/%d", createdID)
		req, _ := http.NewRequest("GET", url, nil)
		rr := executeRequest(mux, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response SuccessResponse[entity.Equipment]
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, createdID, response.Data.ID)
	})

	t.Run("Equipment_Get_NotFound", func(t *testing.T) {
		url := "/equipment/9999" // Несуществующий ID
		req, _ := http.NewRequest("GET", url, nil)
		rr := executeRequest(mux, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Equipment_Create_ValidationError", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/equipment", bytes.NewBufferString(`{"name":"   ","description":"x"}`))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(mux, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var response ErrorResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "validation_error", response.Error.Code)
		assert.Contains(t, response.Error.Message, "name")
	})

	// 3. PUT /equipment/{id} (Update)
	t.Run("Equipment_Update_Success", func(t *testing.T) {
		payload := readTestData(t, "update_equipment.json")
		url := fmt.Sprintf("/equipment/%d", createdID)
		req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(mux, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response SuccessResponse[entity.Equipment]
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, createdID, response.Data.ID)
	})

	// FIXME: тесты сейчас не работают.
	// 4. GET /equipment (List)
	t.Run("Equipment_List_Success", func(t *testing.T) {
		reqList, _ := http.NewRequest("GET", "/equipment", nil)
		rr := executeRequest(mux, reqList)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response SuccessResponse[[]entity.Equipment]
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Data, 1)
	})

	// 5. DELETE /equipment/{id} (Remove)
	t.Run("Equipment_Remove_Success", func(t *testing.T) {
		url := fmt.Sprintf("/equipment/%d", createdID)
		req, _ := http.NewRequest("DELETE", url, nil)
		rr := executeRequest(mux, req)

		// Ваша реализация Delete в handler_equipment.go возвращает 500, если store.Remove()
		// возвращает ошибку (например, NotFoundErr), что не совсем соответствует RESTful стандартам (204 No Content).
		// Однако, если удаление успешно, ожидаем 204.
		assert.Equal(t, http.StatusNoContent, rr.Code)

		// Проверяем, что запись действительно удалена
		reqGet, _ := http.NewRequest("GET", url, nil)
		rrGet := executeRequest(mux, reqGet)
		assert.Equal(t, http.StatusNotFound, rrGet.Code)
	})
	t.Run("Equipment_Export_CSV_Integration", func(t *testing.T) {
		// 1. Создаем запрос с параметром ?format=csv
		req, _ := http.NewRequest("GET", "/equipment?format=csv", nil)
		rr := executeRequest(mux, req)

		// 2. Проверяем статус и заголовки
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "text/csv; charset=utf-8", rr.Header().Get("Content-Type"))
		assert.Contains(t, rr.Header().Get("Content-Disposition"), "attachment")

		// 3. Проверяем содержимое (наличие заголовков, которые мы писали в Store)
		body := rr.Body.String()
		assert.Contains(t, body, "ID")
		assert.Contains(t, body, "Название")

		// Проверка на BOM (чтобы Excel не ругался)
		// Первые три байта должны быть EF BB BF
		assert.True(t, bytes.HasPrefix(rr.Body.Bytes(), []byte("\xEF\xBB\xBF")))
	})
}

// ====================================================================================================
//                             ТЕСТЫ ДЛЯ EXPERIMENT (CRUD & M2M)
// ====================================================================================================

func TestExperimentHandlerCRUD_Integration(t *testing.T) {
	db, mux := setupTestDB(t)
	defer db.Close()

	var expID int // ID для эксперимента
	var eqID1 int // ID для первого оборудования
	var eqID2 int // ID для второго оборудования

	// 1. SETUP: Создание Эксперимента и Оборудования для M2M тестов
	t.Run("Setup_Entities", func(t *testing.T) {
		// 1.1 Создание Эксперимента (POST /experiment)
		expPayload := readTestData(t, "new_experiment.json")
		reqExp, _ := http.NewRequest("POST", "/experiment/", bytes.NewBuffer(expPayload))
		reqExp.Header.Set("Content-Type", "application/json")
		rrExp := executeRequest(mux, reqExp)
		assert.Equal(t, http.StatusCreated, rrExp.Code)
		var experimentResponse SuccessResponse[entity.Experiment]
		err := json.Unmarshal(rrExp.Body.Bytes(), &experimentResponse)
		assert.NoError(t, err)
		expID = experimentResponse.Data.ID

		// 1.2 Создание Оборудования 1 (POST /equipment)
		eqPayload1 := readTestData(t, "new_equipment.json")
		reqEq1, _ := http.NewRequest("POST", "/equipment", bytes.NewBuffer(eqPayload1))
		reqEq1.Header.Set("Content-Type", "application/json")
		rrEq1 := executeRequest(mux, reqEq1)
		assert.Equal(t, http.StatusCreated, rrEq1.Code)
		var equipmentResponse1 SuccessResponse[entity.Equipment]
		err = json.Unmarshal(rrEq1.Body.Bytes(), &equipmentResponse1)
		assert.NoError(t, err)
		eqID1 = equipmentResponse1.Data.ID

		// 1.3 Создание Оборудования 2
		eqPayload2 := []byte(`{"name": "Laser", "description": "High-power beam"}`)
		reqEq2, _ := http.NewRequest("POST", "/equipment", bytes.NewBuffer(eqPayload2))
		reqEq2.Header.Set("Content-Type", "application/json")
		rrEq2 := executeRequest(mux, reqEq2)
		assert.Equal(t, http.StatusCreated, rrEq2.Code)
		var equipmentResponse2 SuccessResponse[entity.Equipment]
		err = json.Unmarshal(rrEq2.Body.Bytes(), &equipmentResponse2)
		assert.NoError(t, err)
		eqID2 = equipmentResponse2.Data.ID
	})

	t.Run("Experiment_Create_ValidationError", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/experiments", bytes.NewBufferString(`{"name":"   ","description":"invalid"}`))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(mux, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var response ErrorResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "validation_error", response.Error.Code)
		assert.Contains(t, response.Error.Message, "name")
	})

	// 2. M2M: POST /experiment/{id}/equipment (Add Equipment)
	t.Run("M2M_AddEquipment_Success", func(t *testing.T) {
		url := fmt.Sprintf("/experiment/%d/equipment", expID)

		// 2.1 Добавление Оборудования 1
		payload1 := fmt.Sprintf(`{"equipment_id": %d}`, eqID1)
		req1, _ := http.NewRequest("POST", url, bytes.NewBufferString(payload1))
		req1.Header.Set("Content-Type", "application/json")
		rr1 := executeRequest(mux, req1)
		assert.Equal(t, http.StatusCreated, rr1.Code, "Добавление первого оборудования должно быть успешным")

		var linkResponse1 SuccessResponse[map[string]int]
		err := json.Unmarshal(rr1.Body.Bytes(), &linkResponse1)
		assert.NoError(t, err)
		assert.Equal(t, expID, linkResponse1.Data["experiment_id"])
		assert.Equal(t, eqID1, linkResponse1.Data["equipment_id"])

		// 2.2 Добавление Оборудования 2
		payload2 := fmt.Sprintf(`{"equipment_id": %d}`, eqID2)
		req2, _ := http.NewRequest("POST", url, bytes.NewBufferString(payload2))
		req2.Header.Set("Content-Type", "application/json")
		rr2 := executeRequest(mux, req2)
		assert.Equal(t, http.StatusCreated, rr2.Code, "Добавление второго оборудования должно быть успешным")

		var linkResponse2 SuccessResponse[map[string]int]
		err = json.Unmarshal(rr2.Body.Bytes(), &linkResponse2)
		assert.NoError(t, err)
		assert.Equal(t, expID, linkResponse2.Data["experiment_id"])
		assert.Equal(t, eqID2, linkResponse2.Data["equipment_id"])
	})

	t.Run("M2M_AddEquipment_BadRequest", func(t *testing.T) {
		url := fmt.Sprintf("/experiment/%d/equipment", expID)
		req, _ := http.NewRequest("POST", url, bytes.NewBufferString(`{"equipment_id":0}`))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(mux, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var response ErrorResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "validation_error", response.Error.Code)
		assert.Contains(t, response.Error.Message, "equipment_id")
	})

	// 3. M2M: GET /experiment/{id}/equipment (List Equipment)
	t.Run("M2M_ListEquipment_Success", func(t *testing.T) {
		url := fmt.Sprintf("/experiment/%d/equipment", expID)
		req, _ := http.NewRequest("GET", url, nil)
		rr := executeRequest(mux, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response SuccessResponse[[]entity.Equipment]
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Data, 2, "Должно быть два прикрепленных оборудования")
	})

	// 4. M2M: DELETE /experiment/{id}/equipment/{equipment_id} (Remove Equipment)
	t.Run("M2M_RemoveEquipment_Success", func(t *testing.T) {
		url := fmt.Sprintf("/experiment/%d/equipment/%d", expID, eqID1)
		req, _ := http.NewRequest("DELETE", url, nil)
		rr := executeRequest(mux, req)
		// handler_experiment.go:RemoveEquipment возвращает 204 No Content
		assert.Equal(t, http.StatusNoContent, rr.Code, "Удаление связи должно быть успешным")

		// Проверка, что осталось только одно оборудование
		listURL := fmt.Sprintf("/experiment/%d/equipment", expID)
		reqList, _ := http.NewRequest("GET", listURL, nil)
		rrList := executeRequest(mux, reqList)

		var response SuccessResponse[[]entity.Equipment]
		err := json.Unmarshal(rrList.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Data, 1, "После удаления должен остаться только один элемент")
	})

	t.Run("M2M_RemoveEquipment_NotFound", func(t *testing.T) {
		// Попытка удалить уже удаленную связь (eqID1)
		url := fmt.Sprintf("/experiment/%d/equipment/%d", expID, eqID1)
		req, _ := http.NewRequest("DELETE", url, nil)
		rr := executeRequest(mux, req)
		assert.Equal(t, http.StatusNotFound, rr.Code, "Попытка удалить несуществующую связь должна вернуть 404")
	})

	// 5. Experiment CRUD (Проверка Get и Delete самого эксперимента)
	t.Run("Experiment_Get_Success", func(t *testing.T) {
		url := fmt.Sprintf("/experiment/%d", expID)
		req, _ := http.NewRequest("GET", url, nil)
		rr := executeRequest(mux, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var response SuccessResponse[entity.Experiment]
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expID, response.Data.ID)
	})

	t.Run("Experiment_GetEquipment_Success", func(t *testing.T) {
		url := fmt.Sprintf("/experiment/%d/equipment/%d", expID, eqID2)
		req, _ := http.NewRequest("GET", url, nil)
		rr := executeRequest(mux, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response SuccessResponse[entity.Equipment]
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, eqID2, response.Data.ID)
	})

	t.Run("Experiment_Remove_Success_Cascades", func(t *testing.T) {
		// Удаляем сам эксперимент
		deleteURL := fmt.Sprintf("/experiment/%d", expID)
		reqDelete, _ := http.NewRequest("DELETE", deleteURL, nil)
		rrDelete := executeRequest(mux, reqDelete)
		assert.Equal(t, http.StatusNoContent, rrDelete.Code)

		// Проверяем, что эксперимент удален
		reqGet, _ := http.NewRequest("GET", deleteURL, nil)
		rrGet := executeRequest(mux, reqGet)
		assert.Equal(t, http.StatusNotFound, rrGet.Code)

		// NOTE: Каскадное удаление связи (ExperimentEquipment) не проверяется
		// напрямую через API, но оно должно произойти на уровне БД.
	})
}

func TestExperimentPluralRoutes_Integration(t *testing.T) {
	db, mux := setupTestDB(t)
	defer db.Close()

	payload := readTestData(t, "new_experiment.json")
	reqCreate, _ := http.NewRequest("POST", "/experiments", bytes.NewBuffer(payload))
	reqCreate.Header.Set("Content-Type", "application/json")
	rrCreate := executeRequest(mux, reqCreate)
	assert.Equal(t, http.StatusCreated, rrCreate.Code)

	var response SuccessResponse[entity.Experiment]
	err := json.Unmarshal(rrCreate.Body.Bytes(), &response)
	assert.NoError(t, err)

	reqGet, _ := http.NewRequest("GET", fmt.Sprintf("/experiments/%d", response.Data.ID), nil)
	rrGet := executeRequest(mux, reqGet)
	assert.Equal(t, http.StatusOK, rrGet.Code)
}

func TestMassiveExperimentEquipmentExport_Integration(t *testing.T) {
	db, mux := setupTestDB(t)
	defer db.Close()

	// 1. Создаем один эксперимент
	expPayload := []byte(`{"name": "Big Data Exp", "description": "Testing 10 items"}`)
	reqExp, _ := http.NewRequest("POST", "/experiment/", bytes.NewBuffer(expPayload))
	reqExp.Header.Set("Content-Type", "application/json")
	rrExp := executeRequest(mux, reqExp)

	var experimentResponse SuccessResponse[entity.Experiment]
	err := json.Unmarshal(rrExp.Body.Bytes(), &experimentResponse)
	assert.NoError(t, err)
	expID := experimentResponse.Data.ID

	// 2. Создаем 10 единиц оборудования и сразу привязываем их к эксперименту
	for i := 1; i <= 10; i++ {
		// Создаем оборудование
		name := fmt.Sprintf("Sensor #%d", i)
		eqPayload := []byte(fmt.Sprintf(`{"name": "%s", "description": "Autogenerated"}`, name))
		reqEq, _ := http.NewRequest("POST", "/equipment", bytes.NewBuffer(eqPayload))
		reqEq.Header.Set("Content-Type", "application/json")
		rrEq := executeRequest(mux, reqEq)

		var equipmentResponse SuccessResponse[entity.Equipment]
		err := json.Unmarshal(rrEq.Body.Bytes(), &equipmentResponse)
		assert.NoError(t, err)

		// Привязываем к эксперименту
		linkPayload := []byte(fmt.Sprintf(`{"equipment_id": %d}`, equipmentResponse.Data.ID))
		linkURL := fmt.Sprintf("/experiment/%d/equipment", expID)
		reqLink, _ := http.NewRequest("POST", linkURL, bytes.NewBuffer(linkPayload))
		reqLink.Header.Set("Content-Type", "application/json")
		executeRequest(mux, reqLink)
	}

	// 3. ТЕСТ ЭКСПОРТА: Запрашиваем CSV
	exportURL := fmt.Sprintf("/experiment/%d/equipment?format=csv", expID)
	reqExport, _ := http.NewRequest("GET", exportURL, nil)
	rrExport := executeRequest(mux, reqExport)

	// Проверяем статус и заголовки
	assert.Equal(t, http.StatusOK, rrExport.Code)
	assert.Equal(t, "text/csv; charset=utf-8", rrExport.Header().Get("Content-Type"))

	// 4. Проверяем содержимое файла
	body := rrExport.Body.String()

	// Проверяем наличие заголовков
	assert.Contains(t, body, "ID Оборудования")
	assert.Contains(t, body, "Название")

	// Проверяем, что в файле реально 10 строк с оборудованием
	// (считаем количество упоминаний "Sensor #")
	for i := 1; i <= 10; i++ {
		assert.Contains(t, body, fmt.Sprintf("Sensor #%d", i))
	}

	// Проверяем наличие BOM (для Excel)
	assert.True(t, bytes.HasPrefix(rrExport.Body.Bytes(), []byte("\xEF\xBB\xBF")))
}
