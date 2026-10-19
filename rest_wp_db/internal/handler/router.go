package handler

import "net/http"

type TaskHandler interface {
	CreateTask(w http.ResponseWriter, r *http.Request)
	GetTasks(w http.ResponseWriter, r *http.Request)
	GetTaskByID(w http.ResponseWriter, r *http.Request)
	CancelTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
}

func RegisterRoutes(mux *http.ServeMux, taskHandler TaskHandler) {
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks", taskHandler.GetTasks)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTaskByID)
	mux.HandleFunc("POST /tasks/{id}/cancel", taskHandler.CancelTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)
}
