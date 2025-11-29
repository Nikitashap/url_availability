package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"url_availability/internal/models"
	"url_availability/internal/pdf"
	"url_availability/internal/queue"
	"url_availability/internal/storage"
)

type API struct {
	Queue *queue.TaskQueue
	Store *storage.Storage
	Ctx   context.Context
}

// Обработка запроса на проверку указанных страниц

func (a *API) CheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ONLY POST METHOD", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Links []string `json:"links"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Unable to decode json", http.StatusBadRequest)
		return
	}
	if len(req.Links) == 0 {
		http.Error(w, "No links found", http.StatusBadRequest)
		return
	}
	task := &models.LinkTask{
		TaskNum: a.Store.NextTaskNum(),
		Links:   req.Links,
		Done:    false,
	}
	a.Store.SaveTask(task)

	select {
	case <-a.Ctx.Done():
		log.Println("Server is shutting down")
	case a.Queue.Ch <- task:
		log.Println("Task queued:", task.TaskNum)
	default:
		log.Println("Queue is full, task saved:", task.TaskNum)
	}

	resp := map[string]any{
		"links_num": task.TaskNum,
		"status":    "task accepted",
	}
	json.NewEncoder(w).Encode(resp)
}

// Обработка запроса на состовление pdf отчета

func (a *API) ReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ONLY POST METHOD", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		List []int `json:"links_list"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Unable to decode json", http.StatusBadRequest)
		return
	}
	tasks := a.Store.GetTasks(req.List)
	data, err := pdf.Generate(tasks)
	if err != nil {
		http.Error(w, "PDF generation error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Write(data)
}
