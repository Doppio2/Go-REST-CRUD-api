package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go_rest_crud/internal/entity"
	"go_rest_crud/internal/repo"
)

// Регулярные выражения для обращения к страницам с определенным оборудованием и к техникe с этим оборудованием.
var (
	ExperimentRe                = regexp.MustCompile(`^/experiments?/?$`)
	ExperimentReWithID          = regexp.MustCompile(`^/experiments?/([0-9]+)$`)
	ExperimentEquipmentRe       = regexp.MustCompile(`^/experiments?/([0-9]+)/equipment/?$`)
	ExperimentEquipmentReWithID = regexp.MustCompile(`^/experiments?/([0-9]+)/equipment/([0-9]+)$`)
)

// Ручка для сущности Equipment.
type ExperimentHandler struct {
	// TODO: Слишком большие названия. Не лучше ли сделать сокращения EStore, EXStore и EEStore????
	ExperimentStore          repo.ExperimentStore
	EquipmentStore           repo.EquipmentStore
	ExperimentEquipmentStore repo.ExperimentEquipmentStore
}

func GetExperimentID(w http.ResponseWriter, r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	var err error
	if len(parts) < 3 {
		log.Printf("ERROR: [GetExperimentID] path too short: %s", r.URL.Path)
		BadRequestHandler(w, "invalid URL")
		return 0, err
	}

	// Получаем id эксперемента из url, куда мы будет добавлять оборудование.
	result, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("ERROR: [GetExperimentID] cannot parse ID '%s' from path: %v", parts[2], err)
		BadRequestHandler(w, "invalid experiment ID")
		return 0, err
	}

	return result, err
}

func GetEquipmentID(w http.ResponseWriter, r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	var err error
	if len(parts) < 5 {
		log.Printf("ERROR: [GetEquipmentID] path too short, expected at least 5 parts: %s", r.URL.Path)
		BadRequestHandler(w, "invalid URL")
		return 0, err
	}

	// Получаем id эксперемента из url, куда мы будет добавлять оборудование.
	result, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Printf("ERROR: [GetEquipmentID] cannot parse equipment ID '%s' from path: %v", parts[4], err)
		BadRequestHandler(w, "invalid equipment ID")
		return 0, err
	}

	return result, err
}

// Конструктор для ручки Experiment.
func NewExperimentHandler(experimentStore repo.ExperimentStore,
	equipmentStore repo.EquipmentStore,
	experimentEquipmentStore repo.ExperimentEquipmentStore,
) *ExperimentHandler {
	return &ExperimentHandler{
		ExperimentStore:          experimentStore,
		EquipmentStore:           equipmentStore,
		ExperimentEquipmentStore: experimentEquipmentStore,
	}
}

// TODO: Она ничем не отличается от функции для equipment. Ее можно объединить в одну.
// Я думаю, что можно сделать интерфейс и сделать какие-то общие функции для взаимодействия.
// Этим я займусь позже, когда все функции реализую.
func (h *ExperimentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var experiment entity.Experiment
	err := json.NewDecoder(r.Body).Decode(&experiment)
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.Create] failed to decode JSON: %v", err)
		BadRequestHandler(w, "invalid request body")
		return
	}

	experiment.Name, experiment.Description, err = normalizeAndValidatePayload(experiment.Name, experiment.Description)
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.Create] validation failed: %v", err)
		ValidationErrorHandler(w, err.Error())
		return
	}

	experiment.CreationDate = time.Now().UTC().Format(time.RFC3339)

	id, err := h.ExperimentStore.Add(experiment)
	if err != nil {
		// TODO: Pass errors to the InternalServerErorHandler function.
		log.Printf("ERROR: [ExperimentHandler.Create] database error: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	experiment.ID = id

	log.Printf("INFO: [ExperimentHandler.Create] successfully created experiment with ID %d", id)

	writeJSONResponse(w, http.StatusCreated, experiment)
}

func (h *ExperimentHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("format") == "csv" {
		filename := "all_experiments.csv"
		if err := h.ExperimentStore.ExportAllToFile(filename); err != nil {
			log.Printf("ERROR: [ExperimentHandler.List] failed to export experiments to %s: %v", filename, err)
			InternalServerErrorHandler(w, r)
			return
		}
		log.Printf("INFO: [ExperimentHandler.List] successfully exported experiments to %s", filename)
		serveCSV(w, r, filename)
		return
	}

	experimentMap, err := h.ExperimentStore.List()
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.List] database error: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	var experimentList []entity.Experiment
	for _, eq := range experimentMap {
		experimentList = append(experimentList, eq)
	}

	writeJSONResponse(w, http.StatusOK, experimentList)
}

func (h *ExperimentHandler) Get(w http.ResponseWriter, r *http.Request) {
	matches := ExperimentReWithID.FindStringSubmatch(r.URL.Path)

	if len(matches) < 2 {
		log.Printf("ERROR: [ExperimentHandler.Get] failed to parse ID from path: %s", r.URL.Path)
		InternalServerErrorHandler(w, r)
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.Get] invalid ID format '%s': %v", matches[1], err)
		BadRequestHandler(w, "invalid ID format")
		return
	}

	experiment, err := h.ExperimentStore.Get(id)
	if err != nil {
		if err == repo.NotFoundErr {
			log.Printf("INFO: [ExperimentHandler.Get] experiment with ID %d not found", id)
			NotFoundHandler(w, r)
		} else {
			log.Printf("ERROR: [ExperimentHandler.Get] database error for ID %d: %v", id, err)
			InternalServerErrorHandler(w, r)
		}

		return
	}

	log.Printf("DEBUG: [ExperimentHandler.Get] successfully retrieved experiment: %+v", experiment)

	writeJSONResponse(w, http.StatusOK, experiment)
}

func (h *ExperimentHandler) Update(w http.ResponseWriter, r *http.Request) {
	matches := ExperimentReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		log.Printf("ERROR: [ExperimentHandler.Update] missing ID in path: %s", r.URL.Path)
		InternalServerErrorHandler(w, r)
		return
	}

	var experiment entity.Experiment

	err := json.NewDecoder(r.Body).Decode(&experiment)
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.Update] failed to decode JSON: %v", err)
		BadRequestHandler(w, "invalid request body")
		return
	}

	experiment.Name, experiment.Description, err = normalizeAndValidatePayload(experiment.Name, experiment.Description)
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.Update] validation failed: %v", err)
		ValidationErrorHandler(w, err.Error())
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.Update] invalid ID format '%s': %v", matches[1], err)
		BadRequestHandler(w, "invalid ID format")
		return
	}

	if err := h.ExperimentStore.Update(id, experiment); err != nil {
		if err == repo.NotFoundErr {
			log.Printf("INFO: [ExperimentHandler.Update] experiment with ID %d not found for update", id)
			NotFoundHandler(w, r)
			return
		}
		log.Printf("ERROR: [ExperimentHandler.Update] database error for ID %d: %v", id, err)
		InternalServerErrorHandler(w, r)
		return
	}

	updatedExperiment, err := h.ExperimentStore.Get(id)
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.Update] failed to fetch updated experiment ID %d: %v", id, err)
		InternalServerErrorHandler(w, r)
		return
	}

	log.Printf("INFO: [ExperimentHandler.Update] successfully updated experiment ID %d", id)
	writeJSONResponse(w, http.StatusOK, updatedExperiment)
}

func (h *ExperimentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	matches := ExperimentReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		log.Printf("ERROR: [ExperimentHandler.Delete] missing ID in path: %s", r.URL.Path)
		InternalServerErrorHandler(w, r)
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.Delete] invalid ID format '%s': %v", matches[1], err)
		BadRequestHandler(w, "invalid ID format")
		return
	}

	if err := h.ExperimentStore.Remove(id); err != nil {
		if err == repo.NotFoundErr {
			log.Printf("INFO: [ExperimentHandler.Delete] experiment with ID %d not found", id)
			NotFoundHandler(w, r)
			return
		}

		log.Printf("ERROR: [ExperimentHandler.Delete] database error for ID %d: %v", id, err)
		InternalServerErrorHandler(w, r)
		return
	}

	log.Printf("INFO: [ExperimentHandler.Delete] successfully removed experiment ID %d", id)
	w.WriteHeader(http.StatusNoContent)
}

// Добавление оборудование к эксперименту.
func (h *ExperimentHandler) AddEquipment(w http.ResponseWriter, r *http.Request) {
	experimentID, err := GetExperimentID(w, r)
	if err != nil {
		return
	}

	if _, err := h.ExperimentStore.Get(experimentID); err != nil {
		if err == repo.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	// Ищем оборудование по полученному из запроса id.
	var payload struct {
		EquipmentID int `json:"equipment_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.AddEquipment] failed to decode JSON payload: %v", err)
		BadRequestHandler(w, "invalid request body")
		return
	}

	if payload.EquipmentID <= 0 {
		ValidationErrorHandler(w, "equipment_id must be a positive integer")
		return
	}

	equipment, err := h.EquipmentStore.Get(payload.EquipmentID)
	if err != nil {
		log.Printf("INFO: [ExperimentHandler.AddEquipment] equipment %d not found: %v", payload.EquipmentID, err)
		NotFoundHandler(w, r)
		return
	}

	// Устанавливаем связь между экспериментом и оборудованием.
	err = h.ExperimentEquipmentStore.Add(experimentID, equipment.ID)
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.AddEquipment] failed to link experiment %d and equipment %d: %v", experimentID, payload.EquipmentID, err)
		InternalServerErrorHandler(w, r)
		return
	}

	// TODO: возможно стоит лучше поставить http.StatusCreated.
	log.Printf("INFO: [ExperimentHandler.AddEquipment] successfully linked equipment %d to experiment %d", payload.EquipmentID, experimentID)
	writeJSONResponse(w, http.StatusCreated, map[string]int{
		"experiment_id": experimentID,
		"equipment_id":  payload.EquipmentID,
	})
}

// Получения списка всего оборудования, которое используется в эксперименте.
func (h *ExperimentHandler) ListEquipment(w http.ResponseWriter, r *http.Request) {
	experimentID, err := GetExperimentID(w, r)
	if err != nil {
		return
	}

	if r.URL.Query().Get("format") == "csv" {
		// experimentID мы уже достали из URL ранее в хэндлере
		filename := fmt.Sprintf("experiment_%d_equipment.csv", experimentID)

		// ВАЖНО: здесь метод принимает ID
		if err := h.ExperimentEquipmentStore.ExportEquipmentToFile(experimentID, filename); err != nil {
			log.Printf("ERROR: [ExperimentHandler.ListEquipment] failed to export CSV for experiment %d: %v", experimentID, err)
			InternalServerErrorHandler(w, r)
			return
		}
		log.Printf("INFO: [ExperimentHandler.ListEquipment] successfully exported CSV for experiment %d", experimentID)
		serveCSV(w, r, filename)
		return
	}

	equipmentMap, err := h.ExperimentEquipmentStore.ListEquipment(experimentID)
	if err != nil {
		log.Printf("ERROR: [ExperimentHandler.ListEquipment] database error for experiment %d: %v", experimentID, err)
		InternalServerErrorHandler(w, r)
		return
	}

	var equipmentList []entity.Equipment
	for _, eq := range equipmentMap {
		equipmentList = append(equipmentList, eq)
	}

	log.Printf("INFO: [ExperimentHandler.ListEquipment] returned %d items for experiment %d", len(equipmentList), experimentID)

	writeJSONResponse(w, http.StatusOK, equipmentList)
}

func (h *ExperimentHandler) GetEquipment(w http.ResponseWriter, r *http.Request) {
	// TODO: нужно сделать функция типа ParserIDsFromUrl(), которая будет всегда возвращать два значения, experimentId и EquipmentId.
	experimentID, err := GetExperimentID(w, r)
	if err != nil {
		return
	}
	equipmentID, err := GetEquipmentID(w, r)
	if err != nil {
		return
	}

	/*
		equipment, err := h.ExperimentEquipmentStore.GetEquipment(experimentID, equipmentID)
		if err != nil {
			NotFoundHandler(w, r)
			return
		}
	*/
	equipment, err := h.ExperimentEquipmentStore.GetEquipment(experimentID, equipmentID)
	if err != nil {
		if err == repo.NotFoundErr {
			log.Printf("INFO: [ExperimentHandler.GetEquipment] equipment %d not found in experiment %d", equipmentID, experimentID)
			NotFoundHandler(w, r)
		} else {
			log.Printf("ERROR: [ExperimentHandler.GetEquipment] database error (Exp: %d, Eq: %d): %v", experimentID, equipmentID, err)
			InternalServerErrorHandler(w, r)
		}
		return
	}

	log.Printf("INFO: [ExperimentHandler.GetEquipment] successfully retrieved equipment %d for experiment %d", equipmentID, experimentID)

	writeJSONResponse(w, http.StatusOK, equipment)
}

// Удаление оборудования из эксперимента
func (h *ExperimentHandler) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	// TODO: нужно сделать функция типа ParserIDsFromUrl(), которая будет всегда возвращать два значения, experimentId и EquipmentId.
	experimentID, err := GetExperimentID(w, r)
	if err != nil {
		return
	}
	equipmentID, err := GetEquipmentID(w, r)
	if err != nil {
		return
	}

	err = h.ExperimentEquipmentStore.Remove(experimentID, equipmentID)

	if err != nil {
		if err == repo.NotFoundErr {
			// Логируем попытку удаления несуществующей связи
			log.Printf("INFO: [ExperimentHandler.DeleteEquipment] link not found for Exp: %d, Eq: %d", experimentID, equipmentID)
			NotFoundHandler(w, r)
			return
		}
		log.Printf("ERROR: [ExperimentHandler.DeleteEquipment] database error (Exp: %d, Eq: %d): %v", experimentID, equipmentID, err)
		InternalServerErrorHandler(w, r)
		return
	}

	log.Printf("INFO: [ExperimentHandler.DeleteEquipment] successfully unlinked equipment %d from experiment %d", equipmentID, experimentID)
	w.WriteHeader(http.StatusNoContent)
}

// Функция для получения списка всех экспериментов, где используется это оборудование.
// NOTE: нигде не использутся.
//func (h *ExperimentHandler) ListExperiments(w http.ResponseWriter, r *http.Request) {
//}

func (h *ExperimentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	// -- Операции, связанные с экспериментами.
	case r.Method == http.MethodPost && ExperimentRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	case r.Method == http.MethodGet && ExperimentRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && ExperimentReWithID.MatchString(r.URL.Path):
		h.Get(w, r)
		return
	case r.Method == http.MethodPut && ExperimentReWithID.MatchString(r.URL.Path):
		h.Update(w, r)
		return
	case r.Method == http.MethodDelete && ExperimentReWithID.MatchString(r.URL.Path):
		h.Delete(w, r)
		return
		// -- Операции, связанные с экипировкой, которая принадлежит экспериментам.
	case r.Method == http.MethodPost && ExperimentEquipmentRe.MatchString(r.URL.Path):
		h.AddEquipment(w, r)
		return
	case r.Method == http.MethodGet && ExperimentEquipmentRe.MatchString(r.URL.Path):
		h.ListEquipment(w, r)
		return
	case r.Method == http.MethodGet && ExperimentEquipmentReWithID.MatchString(r.URL.Path):
		h.GetEquipment(w, r)
		return
	case r.Method == http.MethodDelete && ExperimentEquipmentReWithID.MatchString(r.URL.Path):
		h.DeleteEquipment(w, r)
		return
	default:
		http.NotFound(w, r)
		return
	}
}
