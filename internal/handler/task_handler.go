package handler

import (
	"encoding/json"
	"net/http"

	"todo-list/internal/model"
	"todo-list/internal/repository"
	"todo-list/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{service: svc}
}

func (h *TaskHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/tasks", h.GetAllTasks).Methods("GET")
	router.HandleFunc("/api/tasks", h.CreateTask).Methods("POST")
	router.HandleFunc("/api/tasks/{id}", h.GetTaskByID).Methods("GET")
	router.HandleFunc("/api/tasks/{id}", h.UpdateTask).Methods("PUT")
	router.HandleFunc("/api/tasks/{id}", h.DeleteTask).Methods("DELETE")
	router.HandleFunc("/api/tasks/{id}/status", h.UpdateTaskStatus).Methods("PATCH")
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	filters := repository.TaskFilters{}

	if status := r.URL.Query().Get("status"); status != "" {
		s := model.Status(status)
		filters.Status = &s
	}
	if priority := r.URL.Query().Get("priority"); priority != "" {
		p := model.Priority(priority)
		filters.Priority = &p
	}
	if dueDate := r.URL.Query().Get("due_date"); dueDate != "" {
		filters.DueDate = &dueDate
	}

	tasks, err := h.service.GetAll(r.Context(), filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	task, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req service.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.service.Create(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	var req service.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		if err.Error() == "task not found" {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if err.Error() == "task not found" {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	var req struct {
		Status model.Status `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.service.UpdateStatus(r.Context(), id, req.Status)
	if err != nil {
		if err.Error() == "task not found" {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
