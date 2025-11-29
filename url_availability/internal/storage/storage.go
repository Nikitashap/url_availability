package storage

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"url_availability/internal/models"
)

type Storage struct {
	mu        sync.Mutex
	tasks     map[int]*models.LinkTask
	counter   int
	taskFile  string
	countFile string
}

func NewStorage(taskFile, countFile string) *Storage {
	s := &Storage{
		tasks:     make(map[int]*models.LinkTask),
		taskFile:  taskFile,
		countFile: countFile,
	}

	s.load()
	return s
}

func (s *Storage) load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if data, err := os.ReadFile(s.taskFile); err == nil {
		json.Unmarshal(data, &s.tasks)
	}

	if data, err := os.ReadFile(s.countFile); err == nil {
		s.counter, _ = strconv.Atoi(string(data))
	}
}

func (s *Storage) NextTaskNum() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	os.WriteFile(s.countFile, []byte(strconv.Itoa(s.counter)), 0644)
	return s.counter
}

func (s *Storage) SaveTask(task *models.LinkTask) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.TaskNum] = task
	s.flush()
}

func (s *Storage) GetTasks(nums []int) []*models.LinkTask {
	s.mu.Lock()
	defer s.mu.Unlock()

	var res []*models.LinkTask
	for _, n := range nums {
		if t, ok := s.tasks[n]; ok {
			res = append(res, t)
		}
	}
	return res
}

func (s *Storage) AllUndone() []*models.LinkTask {
	s.mu.Lock()
	defer s.mu.Unlock()

	var res []*models.LinkTask
	for _, t := range s.tasks {
		if !t.Done {
			res = append(res, t)
		}
	}
	return res
}

func (s *Storage) flush() {
	data, _ := json.MarshalIndent(s.tasks, "", "  ")
	os.WriteFile(s.taskFile, data, 0644)
}
