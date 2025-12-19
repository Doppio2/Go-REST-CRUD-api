package handler

import (
	"time"
	//"fmt"
	"net/http"
	"regexp"
	"encoding/json"
	"log"
	"strconv"

	"go_rest_crud/internal/repo"
	"go_rest_crud/internal/entity"
)

// Регулярные выражения для обращения к страницам с определенным оборудованием.
var (
	EquipmentRe = regexp.MustCompile(`^/equipment/?$`)
	EquipmentReWithID = regexp.MustCompile(`^/equipment/([0-9]+)$`)
)

// Ручка для сущности Equipment.
type EquipmentHandler struct {
	store repo.EquipmentStore
}

// Конструктор для ручки Equipment.
func NewEquipmentHandler(s repo.EquipmentStore) *EquipmentHandler {
	return &EquipmentHandler {
		store: s,
	}
}

// Функции обработчики запросов.
func (h *EquipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	// NOTE: rename to just "e".
	var equipment entity.Equipment
	err := json.NewDecoder(r.Body).Decode(&equipment)
	equipment.CreationDate = time.Now().UTC().Format(time.RFC3339)

	if err != nil {
		// TODO: Log later.
		// TODO: Pass errors to the InternalServerErorHandler function.
		log.Fatal("Cant get json body: ", err)
		InternalServerErrorHandler(w, r)
		return 
	}

	id, err := h.store.Add(equipment)
	if err != nil {
		// TODO: Log later.
		// TODO: Pass errors to the InternalServerErorHandler function.
		log.Fatal("Can not add equipment to the database", err)
		InternalServerErrorHandler(w, r)
		return
	}

	equipment.ID = id

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipment)
}

// Получение всех записей из бд.
func (h *EquipmentHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("format") == "csv" {
		filename := "all_equipment.csv"
		if err := h.store.ExportAllToFile(filename); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		serveCSV(w, r, filename) // Вынес отправку в отдельный метод для чистоты
		return
	}

    equipmentMap, err := h.store.List()
    if err != nil {
		// TODO: Log later.
        InternalServerErrorHandler(w, r)
        return
    }

    var equipmentList []entity.Equipment
    for _, eq := range equipmentMap {
        equipmentList = append(equipmentList, eq)
    }

    jsonBytes, err := json.Marshal(equipmentList)
    if err != nil {
        InternalServerErrorHandler(w, r)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonBytes)
}

// Получение 
func (h *EquipmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	matches := EquipmentReWithID.FindStringSubmatch(r.URL.Path)

	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		// TODO: Log later.
		log.Fatal("Can't get element ID: ", err)
	}

	equipment, err := h.store.Get(id)
	if err != nil {
		if err == repo.NotFoundErr {
			NotFoundHandler(w, r)
		} else {
			InternalServerErrorHandler(w, r)
		}
		
		return
	}

	jsonBytes, err := json.Marshal(equipment)
	if err != nil {
		InternalServerErrorHandler(w, r)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *EquipmentHandler) Update(w http.ResponseWriter, r *http.Request) {
    matches := EquipmentReWithID.FindStringSubmatch(r.URL.Path)
    if len(matches) < 2 {
        InternalServerErrorHandler(w, r)
        return
    }

    var equipment entity.Equipment
	
    err := json.NewDecoder(r.Body).Decode(&equipment)
	if err != nil {
        InternalServerErrorHandler(w, r)
        return
    }

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		// TODO: Log later.
		log.Fatal("Can't get elemetn ID: ", err)
	}

    if err := h.store.Update(id, equipment); err != nil {
        if err == repo.NotFoundErr {
            NotFoundHandler(w, r)
            return
        }
        InternalServerErrorHandler(w, r)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *EquipmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
    matches := EquipmentReWithID.FindStringSubmatch(r.URL.Path)
    if len(matches) < 2 {
        InternalServerErrorHandler(w, r)
        return
    }

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		// TODO: Log later.
		log.Fatal("Can't get element ID: ", err)
	}

    if err := h.store.Remove(id); err != nil {
        InternalServerErrorHandler(w, r)
        return
    }
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
