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

// Регулярные выражения для обращения к страницам с определенным оборудованием и к техники с этим оборудованием.
var (
	// TODO: Я не знаю, работает ли это вообще.
	ExperimentRe                = regexp.MustCompile(`^/experiment/*$`)
	ExperimentReWithID          = regexp.MustCompile(`^/experiment/([a-z0-9]+(?:-[a-z0-9]+)*)$`)
	ExperimentEquipmentRe       = regexp.MustCompile(`^/experiment/*$/equipment/*$`)
	ExperimentEquipmentReWithID = regexp.MustCompile(`^/experiment/([a-z0-9]+(?:-[a-z0-9]+)*)$/equipment/([a-z0-9]+(?:-[a-z0-9]+)*)$`)
)

// Ручка для сущности Equipment.
type ExperimentHandler struct {
	ExperimentStore               repo.ExperimentStore
	ExperimentEquipmentStore repo.ExperimentEquipmentStore
}

// Конструктор для ручки Experiment.
func NewExperimentHandler(experimentStore repo.ExperimentStore, experimentEquipmentStore repo.ExperimentEquipmentStore) *ExperimentHandler {
	return &ExperimentHandler {
		ExperimentStore: experimentStore,
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
		// TODO: Pass errors to the InternalServerErorHandler function.
		log.Fatal("Cant get json body: ", err)
		InternalServerErrorHandler(w, r)
		return 
	}

	resourceID := slug.Make(experiment.Name)
	err = h.ExperimentStore.Add(resourceID, experiment)
	if err != nil {
		// TODO: Pass errors to the InternalServerErorHandler function.
		log.Fatal("Can not add experiment to the database", err)
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ExperimentHandler) List(w http.ResponseWriter, r *http.Request) {
    experimentMap, err := h.ExperimentStore.List()
    if err != nil {
        InternalServerErrorHandler(w, r)
        return
    }

    var experimentList []entity.Experiment
    for _, eq := range experimentMap {
        experimentList = append(experimentList, eq)
    }

    jsonBytes, err := json.Marshal(experimentList)
    if err != nil {
        InternalServerErrorHandler(w, r)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonBytes)
}
func (h *ExperimentHandler) Get(w http.ResponseWriter, r *http.Request) {
	matches := ExperimentReWithID.FindStringSubmatch(r.URL.Path)

	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	experiment, err := h.ExperimentStore.Get(matches[1])
	if err != nil {
		if err == repo.NotFoundErr {
			NotFoundHandler(w, r)
		} else {
			InternalServerErrorHandler(w, r)
		}
		
		return
	}

	jsonBytes, err := json.Marshal(experiment)
	if err != nil {
		InternalServerErrorHandler(w, r)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h *ExperimentHandler) Update(w http.ResponseWriter, r *http.Request) {
    matches := ExperimentReWithID.FindStringSubmatch(r.URL.Path)
    if len(matches) < 2 {
        InternalServerErrorHandler(w, r)
        return
    }

    var experiment entity.Experiment
	
    err := json.NewDecoder(r.Body).Decode(&experiment)
	if err != nil {
        InternalServerErrorHandler(w, r)
        return
    }

    if err := h.ExperimentStore.Update(matches[1], experiment); err != nil {
        if err == repo.NotFoundErr {
            NotFoundHandler(w, r)
            return
        }
        InternalServerErrorHandler(w, r)
        return
    }

    w.WriteHeader(http.StatusOK)

}

func (h *ExperimentHandler) Delete(w http.ResponseWriter, r *http.Request) {
    matches := EquipmentReWithID.FindStringSubmatch(r.URL.Path)
    if len(matches) < 2 {
        InternalServerErrorHandler(w, r)
        return
    }
    if err := h.ExperimentStore.Remove(matches[1]); err != nil {
        InternalServerErrorHandler(w, r)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (h *ExperimentHandler) AddEquipment(w http.ResponseWriter, r *http.Request) {
}

func (h *ExperimentHandler) ListEquipment(w http.ResponseWriter, r *http.Request) {
}

func (h *ExperimentHandler) GetEquipment(w http.ResponseWriter, r *http.Request) {
}

func (h *ExperimentHandler) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
}

func (h *ExperimentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
		// -- Операции, связанные с экспериментами.
	case r.Method == http.MethodPost && ExperimentRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	case r.Method == http.MethodGet && ExperimentRe.MatchString(r.URL.Path):
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
		// TODO: Я думаю, что это можно как-то сократить, но пока я не знаю как.
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
		return
	}
}
