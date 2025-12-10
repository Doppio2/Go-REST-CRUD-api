package handler

import (
	"net/http"
	"regexp"
	"encoding/json"
	"log"
	"strconv"

	"go_rest_crud/internal/repo"
	"go_rest_crud/internal/entity"
)

// Регулярные выражения для обращения к страницам с определенным оборудованием и к техники с этим оборудованием.
var (
	// TODO: пока что временно тут чисто числа в url, но я пока не особо хочу заморачиваться с этим всем. Так что пусть будет так.
	ExperimentRe                = regexp.MustCompile(`^/experiment/?$`)
	ExperimentReWithID          = regexp.MustCompile(`^/experiment/([0-9]+)$`)
	ExperimentEquipmentRe       = regexp.MustCompile(`^/experiment/([0-9]+)/equipment/?$`)
	ExperimentEquipmentReWithID = regexp.MustCompile(`^/experiment/([0-9]+)/equipment/([0-9]+)$`)
)

// Ручка для сущности Equipment.
type ExperimentHandler struct {
	ExperimentStore               repo.ExperimentStore
	ExperimentEquipmentStore      repo.ExperimentEquipmentStore
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

	err = h.ExperimentStore.Add(experiment)
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

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		// TODO: Log later.
		log.Fatal("Can't get element ID: ", err)
	}

	experiment, err := h.ExperimentStore.Get(id)
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

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		// TODO: Log later.
		log.Fatal("Can't get element ID: ", err)
	}

    if err := h.ExperimentStore.Update(id, experiment); err != nil {
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
    matches := ExperimentReWithID.FindStringSubmatch(r.URL.Path)
    if len(matches) < 2 {
        InternalServerErrorHandler(w, r)
        return
    }

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		// TODO: Log later.
		log.Fatal("Can't get element ID: ", err)
	}

    if err := h.ExperimentStore.Remove(id); err != nil {
        InternalServerErrorHandler(w, r)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

// Добавление оборудование к эксперименту. 
func (h *ExperimentHandler) AddEquipment(w http.ResponseWriter, r *http.Request) {
}

// Удаление оборудования из эксперимента
func (h *ExperimentHandler) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
}

// Получения списка всего оборудования, которое используется в эксперименте.
func (h *ExperimentHandler) ListEquipment(w http.ResponseWriter, r *http.Request) {
}

// Функция для получения списка всех экспериментов, где используется это оборудование.
func (h *ExperimentHandler) ListExperiments(w http.ResponseWriter, r *http.Request) {
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
		// TODO: Добавить метод GetEquipment.
		// TODO: И, возможно, добавить метод GetExperimen
		// TODO: Не знаю, нужен ли метод GetEquipment(). Возможно да, возможно нет.
		// -- Операции, связанные с экипировкой, которая принадлежит экспериментам.
	case r.Method == http.MethodPost && ExperimentEquipmentRe.MatchString(r.URL.Path):
		h.AddEquipment(w, r)
		return
	case r.Method == http.MethodGet && ExperimentEquipmentRe.MatchString(r.URL.Path):
		h.ListEquipment(w, r)
		return
		// TODO: Если я захочу реализовать этот метод, то мне лучше положить его в 
		// handler_equipment.go и обрабатывать это там. Но мне кажется, что это слишком излишне.
//	case r.Method == http.MethodGet && ExperimentEquipmentRe.MatchString(r.URL.Path):
//		h.ListExperiments(w, r)
//		return
	case r.Method == http.MethodDelete && ExperimentEquipmentReWithID.MatchString(r.URL.Path):
		h.DeleteEquipment(w, r)
		return
	default:
		return
	}
}
