package worker

import (
	"context"
	"log"
	"sync"

	"url_availability/internal/checker"
	"url_availability/internal/models"
	"url_availability/internal/storage"
)

func StartWorker(
	ctx context.Context,
	wg *sync.WaitGroup,
	tasks <-chan *models.LinkTask,
	st *storage.Storage,
) {
	defer wg.Done()
	log.Println("Worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopping")
			return

		case task, ok := <-tasks:
			if !ok {
				log.Println("Task channel closed")
				return
			}

			log.Println("Worker processing task:", task.TaskNum)

			task.Statuses = make(map[string]models.LinkStatus)
			for _, link := range task.Links {
				task.Statuses[link] = checker.CheckLink(link)
			}

			task.Done = true
			st.SaveTask(task)

			log.Println("Worker finished task:", task.TaskNum)
		}
	}
}
