package handler

import (
	"net/http"
	"regexp"
	"encoding/json"
	"log"

	"go_rest_crud/internal/repo"
	"go_rest_crud/internal/entity"

	"github.com/gosimple/slug"
)

// Регулярные выражения для обращения к страницам с определенным оборудованием.
var (
	EquipmentRe       = regexp.MustCompile(`^/equipment/*$`)
	EquipmentReWithID = regexp.MustCompile(`^/equipment/([a-z0-9]+(?:-[a-z0-9]+)*)$`)
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
// TODO: подумать над тем, нужно ли оставлять это функцией. И зачем вообще в конце названий equipment стоит? 
// Едва ли там какой-то конфликт имен есть? Это касается всех функций, которые идут следом.
// Создание записи в бд.
func (h *EquipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var equipment entity.Equipment

	err := json.NewDecoder(r.Body).Decode(&equipment)

	if err != nil {
		// TODO: Pass errors to the InternalServerErorHandler function.
		log.Fatal("Cant get json body: ", err)
		InternalServerErrorHandler(w, r)
		return 
	}

	resourceID := slug.Make(equipment.Name)
	err = h.store.Add(resourceID, equipment)
	if err != nil {
		// TODO: Pass errors to the InternalServerErorHandler function.
		log.Fatal("Can not add equipment to the database", err)
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Получение всех записей из бд.
func (h *EquipmentHandler) List(w http.ResponseWriter, r *http.Request) {
    equipmentMap, err := h.store.List()
    if err != nil {
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

	equipment, err := h.store.Get(matches[1])
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

    if err := h.store.Update(matches[1], equipment); err != nil {
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
    if err := h.store.Remove(matches[1]); err != nil {
        InternalServerErrorHandler(w, r)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

// Функций для 
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
		return
	}
}
