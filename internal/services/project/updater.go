package projectservice

import (
	"project-service/internal/domain/models"
	"sync"
)

type ProjectUpdater struct {
	updates chan *models.Project
	mu      sync.Mutex
	closed  bool
}

func NewProjectUpdater() *ProjectUpdater {
	return &ProjectUpdater{
		updates: make(chan *models.Project, 100),
		closed:  false,
	}
}

func (pu *ProjectUpdater) Publish(project *models.Project) {
	pu.mu.Lock()
	defer pu.mu.Unlock()
	if pu.closed {
		return
	}
	pu.updates <- project
}

func (pu *ProjectUpdater) Subscribe() <-chan *models.Project {
	return pu.updates
}

func (pu *ProjectUpdater) Close() {
	pu.mu.Lock()
	defer pu.mu.Unlock()
	if !pu.closed {
		close(pu.updates)
		pu.closed = true
	}
}
