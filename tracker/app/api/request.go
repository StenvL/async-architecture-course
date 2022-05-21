package api

import (
	"github.com/StenvL/async-architecture-course/tracker/app/model"
	"github.com/google/uuid"
)

type newTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (ntr newTaskRequest) toModel() model.Task {
	taskUUID, _ := uuid.NewUUID()
	return model.Task{
		PublicID:    taskUUID,
		Title:       ntr.Title,
		Description: ntr.Description,
	}
}
