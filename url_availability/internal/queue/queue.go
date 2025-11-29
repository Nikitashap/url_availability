package queue

import "url_availability/internal/models"

type TaskQueue struct {
	Ch chan *models.LinkTask
}

func NewQueue(size int) *TaskQueue {
	return &TaskQueue{
		Ch: make(chan *models.LinkTask, size),
	}
}
