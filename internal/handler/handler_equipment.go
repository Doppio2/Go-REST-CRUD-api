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

var (
	EquipmentRe       = regexp.MustCompile(`^/equipment/*$`)
	EquipmentReWithID = regexp.MustCompile(`^/equipment/([a-z0-9]+(?:-[a-z0-9]+)*)$`)
)

// TODO: move it to handler_error.go
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("404 Not Found"))
}

// TODO: move it to handler_home.go
type HomeHandler struct{}

func (h *HomeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("This is my home page"))
}

type EquipmentHandler struct {
	//store repo.EquipmentStore
	store repo.EquipmentStore
}

func NewEquipmentHandler(s repo.EquipmentStore) *EquipmentHandler {
	return &EquipmentHandler {
		store: s,
	}
}

//func (h *EquipmentHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
//	writer.Write([]byte("Equipment page"))
//}

func (h *EquipmentHandler) CreateEquipment(w http.ResponseWriter, r *http.Request) {

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

func (h *EquipmentHandler) ListEquipment(w http.ResponseWriter, r *http.Request) {
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

func (h *EquipmentHandler) GetEquipment(w http.ResponseWriter, r *http.Request) {
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

func (h *EquipmentHandler) UpdateEquipment(w http.ResponseWriter, r *http.Request) {
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

func (h *EquipmentHandler) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
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

func (h *EquipmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && EquipmentRe.MatchString(r.URL.Path):
		h.CreateEquipment(w, r)
		return
	case r.Method == http.MethodGet && EquipmentRe.MatchString(r.URL.Path):
		h.ListEquipment(w, r)
		return
	case r.Method == http.MethodGet && EquipmentReWithID.MatchString(r.URL.Path):
		h.GetEquipment(w, r)
		return
	case r.Method == http.MethodPut && EquipmentReWithID.MatchString(r.URL.Path):
		h.UpdateEquipment(w, r)
		return
	case r.Method == http.MethodDelete && EquipmentReWithID.MatchString(r.URL.Path):
		h.DeleteEquipment(w, r)
		return
	default:
		return
	}
}
