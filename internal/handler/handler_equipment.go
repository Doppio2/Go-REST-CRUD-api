package handler

import (
	"time"
	//"fmt"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"go_rest_crud/internal/entity"
	"go_rest_crud/internal/repo"
)

// Регулярные выражения для обращения к страницам с определенным оборудованием.
var (
	EquipmentRe       = regexp.MustCompile(`^/equipment/?$`)
	EquipmentReWithID = regexp.MustCompile(`^/equipment/([0-9]+)$`)
)

// Ручка для сущности Equipment.
type EquipmentHandler struct {
	store repo.EquipmentStore
}

// Конструктор для ручки Equipment.
func NewEquipmentHandler(s repo.EquipmentStore) *EquipmentHandler {
	return &EquipmentHandler{
		store: s,
	}
}

// Функции обработчики запросов.
func (h *EquipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	// NOTE: rename to just "e".
	var equipment entity.Equipment
	err := json.NewDecoder(r.Body).Decode(&equipment)
	if err != nil {
		log.Printf("ERROR: [EquipmentHandler.Create] failed to decode JSON: %v", err)
		BadRequestHandler(w, "invalid request body")
		return
	}

	equipment.Name, equipment.Description, err = normalizeAndValidatePayload(equipment.Name, equipment.Description)
	if err != nil {
		log.Printf("ERROR: [EquipmentHandler.Create] validation failed: %v", err)
		ValidationErrorHandler(w, err.Error())
		return
	}

	equipment.CreationDate = time.Now().UTC().Format(time.RFC3339)

	id, err := h.store.Add(equipment)
	if err != nil {
		// TODO: Pass errors to the InternalServerErorHandler function.
		//log.Fatal("Can not add equipment to the database", err)
		log.Printf("ERROR: [EquipmentHandler.Create] database error: %v", err)

		InternalServerErrorHandler(w, r)
		return
	}

	equipment.ID = id

	writeJSONResponse(w, http.StatusCreated, equipment)
}

// Получение всех записей из бд.
func (h *EquipmentHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("format") == "csv" {
		filename := "all_equipment.csv"
		if err := h.store.ExportAllToFile(filename); err != nil {
			log.Printf("ERROR: [EquipmentHandler.List] failed to export CSV to %s: %v", filename, err)
			InternalServerErrorHandler(w, r)
			return
		}
		serveCSV(w, r, filename) // Вынес отправку в отдельный метод для чистоты
		return
	}

	equipmentMap, err := h.store.List()
	if err != nil {
		log.Printf("ERROR: [EquipmentHandler.List] database error: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	var equipmentList []entity.Equipment
	for _, eq := range equipmentMap {
		equipmentList = append(equipmentList, eq)
	}

	writeJSONResponse(w, http.StatusOK, equipmentList)
}

// Получение
func (h *EquipmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	matches := EquipmentReWithID.FindStringSubmatch(r.URL.Path)

	if len(matches) < 2 {
		log.Printf("ERROR: [EquipmentHandler.Get] failed to extract ID from path: %s", r.URL.Path)
		InternalServerErrorHandler(w, r)
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		// TODO: Log later.
		log.Printf("ERROR: [EquipmentHandler.Get] invalid ID format '%s': %v", matches[1], err)
		BadRequestHandler(w, "invalid ID format")
		return
	}

	equipment, err := h.store.Get(id)
	if err != nil {
		if err == repo.NotFoundErr {
			log.Printf("INFO: [EquipmentHandler.Get] equipment with ID %d not found", id)
			NotFoundHandler(w, r)
		} else {
			log.Printf("ERROR: [EquipmentHandler.Get] database error for ID %d: %v", id, err)
			InternalServerErrorHandler(w, r)
		}

		return
	}

	writeJSONResponse(w, http.StatusOK, equipment)
}

func (h *EquipmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	matches := EquipmentReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		log.Printf("ERROR: [EquipmentHandler.Update] missing ID in path: %s", r.URL.Path)
		InternalServerErrorHandler(w, r)
		return
	}

	var equipment entity.Equipment

	err := json.NewDecoder(r.Body).Decode(&equipment)
	if err != nil {
		log.Printf("ERROR: [EquipmentHandler.Update] failed to decode JSON: %v", err)
		BadRequestHandler(w, "invalid request body")
		return
	}

	equipment.Name, equipment.Description, err = normalizeAndValidatePayload(equipment.Name, equipment.Description)
	if err != nil {
		log.Printf("ERROR: [EquipmentHandler.Update] validation failed: %v", err)
		ValidationErrorHandler(w, err.Error())
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Printf("ERROR: [EquipmentHandler.Update] invalid ID format '%s': %v", matches[1], err)
		BadRequestHandler(w, "invalid ID format")
		return
	}

	if err := h.store.Update(id, equipment); err != nil {
		if err == repo.NotFoundErr {
			log.Printf("INFO: [EquipmentHandler.Update] attempt to update non-existent ID %d", id)
			NotFoundHandler(w, r)
			return
		}
		log.Printf("ERROR: [EquipmentHandler.Update] database error for ID %d: %v", id, err)
		InternalServerErrorHandler(w, r)
		return
	}

	updatedEquipment, err := h.store.Get(id)
	if err != nil {
		log.Printf("ERROR: [EquipmentHandler.Update] failed to fetch updated equipment ID %d: %v", id, err)
		InternalServerErrorHandler(w, r)
		return
	}

	log.Printf("INFO: [EquipmentHandler.Update] successfully updated equipment ID %d", id)
	writeJSONResponse(w, http.StatusOK, updatedEquipment)
}

func (h *EquipmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	matches := EquipmentReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		log.Printf("ERROR: [EquipmentHandler.Delete] missing ID in path: %s", r.URL.Path)
		InternalServerErrorHandler(w, r)
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Printf("ERROR: [EquipmentHandler.Delete] invalid ID format '%s': %v", matches[1], err)
		BadRequestHandler(w, "invalid ID format")
		return
	}

	if err := h.store.Remove(id); err != nil {
		if err == repo.NotFoundErr {
			log.Printf("INFO: [EquipmentHandler.Delete] equipment with ID %d not found", id)
			NotFoundHandler(w, r)
			return
		}

		log.Printf("ERROR: [EquipmentHandler.Delete] database error for ID %d: %v", id, err)
		InternalServerErrorHandler(w, r)
		return
	}

	log.Printf("INFO: [EquipmentHandler.Delete] successfully removed equipment ID %d", id)
	w.WriteHeader(http.StatusNoContent)
}

// Обработка запросов.
func (h *EquipmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && EquipmentRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	case r.Method == http.MethodGet && EquipmentRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && EquipmentReWithID.MatchString(r.URL.Path):
		h.Get(w, r)
		return
	case r.Method == http.MethodPut && EquipmentReWithID.MatchString(r.URL.Path):
		h.Update(w, r)
		return
	case r.Method == http.MethodDelete && EquipmentReWithID.MatchString(r.URL.Path):
		h.Delete(w, r)
		return
	default:
		http.NotFound(w, r)
		return
	}
}
