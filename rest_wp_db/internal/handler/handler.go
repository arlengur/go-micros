package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	errs "github.com/arlengur/go_assessment/internal/errors"
	"github.com/arlengur/go_assessment/internal/service"
)

type handler struct {
	service service.Service
}

func NewHandler(service service.Service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuid, err := h.service.CreateTask(r.Context(), req.Type, string(req.Payload))
	if err != nil {
		log.Printf("Error creating task: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Task created with UUID: %s", uuid)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(uuid); err != nil {
		log.Printf("JSON encoding error: %v", err)
	}
}

func (h *handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetTasks(r.Context())
	if err != nil {
		log.Printf("Error getting tasks: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully retrieved %d tasks", len(tasks))
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(tasks); err != nil {
		log.Printf("JSON encoding erro: %v", err)
	}
}

func (h *handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	task, err := h.service.GetTaskByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Printf("Error getting task by id: %v", err)
		return
	}

	log.Printf("Successfully retrieved task with id: %s", id)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("JSON encoding erro: %v", err)
	}
}

func (h *handler) CancelTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.CancelTask(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Printf("Error canceling task: %v", err)
		return
	}

	log.Printf("Task canceled successfully: %s", id)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Printf("Error deleting task: %v", err)
		return
	}

	log.Printf("Task deleted successfully: %s", id)
	w.WriteHeader(http.StatusOK)
}
