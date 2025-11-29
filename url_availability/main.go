package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"url_availability/internal/handlers"
	"url_availability/internal/queue"
	"url_availability/internal/storage"
	"url_availability/internal/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := storage.NewStorage("data/tasks.json", "data/counter.txt")
	queue := queue.NewQueue(100)

	// Восстановление незавершённых задач при запуске
	for _, task := range store.AllUndone() {
		queue.Ch <- task
	}
	log.Println("Unfinished tasks:", len(store.AllUndone()))

	// Создаем WaitGroup для параллельной обработки запросов
	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go worker.StartWorker(ctx, &wg, queue.Ch, store)
	}

	api := &handlers.API{
		Queue: queue,
		Store: store,
		Ctx:   ctx,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/check", api.CheckHandler)
	mux.HandleFunc("/report", api.ReportHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Описание процесса остановки сервиса
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("Shutting down")
		cancel()
		server.Shutdown(context.Background())
	}()

	log.Println("Server running on :8080")

	server.ListenAndServe()

	wg.Wait()
	log.Println("Server stopped")
}
